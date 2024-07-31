package productrepository

import (
	"fmt"
	"net/http"

	"hometest1/core/entity"
	"hometest1/core/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.ProductRepo {
	return &repo{db}
}

func (r *repo) GetProductBySerials(serials []string) ([]*entity.Product, error) {
	var result []*entity.Product
	err := r.db.Where("serial in (?)", serials).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repo) GetProductByIDs(ids []int64) ([]*entity.Product, error) {
	var result []*entity.Product
	err := r.db.Where("id in (?)", ids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repo) SubmitCheckout(payload *entity.Checkout) (err error) {
	// begin transaction
	tx := r.db.Begin()
	defer func() {
		if rc := recover(); rc != nil {
			tx.Rollback()
			switch x := rc.(type) {
			case string:
				err = entity.NewError(x, http.StatusInternalServerError)
			case error:
				err = entity.NewError(x.Error(), http.StatusInternalServerError)
			default:
				err = entity.NewError("unknown panic", http.StatusInternalServerError)
			}
		}
	}()
	err = tx.Error
	if err != nil {
		return
	}

	// lock for update product quantity
	var mapProdQty map[int64]*entity.ProductQuantity
	productIDs := r.pluckProductIDFromCheckoutItems(payload.Items)
	mapProdQty, err = r.lockAndMapProductQuantity(productIDs, tx)
	if err != nil {
		err = entity.NewError(err.Error(), http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	// validate and update product quantity
	for _, item := range payload.Items {
		newQuantity := mapProdQty[item.Product.ID].Quantity - item.Quantity
		if newQuantity < 0 {
			err = entity.NewError(
				fmt.Sprintf("checkout item %s(%s) exceeds existing quantity, only %d items remaining",
					item.Product.Name, item.Product.Serial, mapProdQty[item.Product.ID].Quantity),
				http.StatusBadRequest)
			tx.Rollback()
			return
		}

		// update table product_quantity
		mapProdQty[item.Product.ID].Quantity = newQuantity
		err = tx.Save(mapProdQty[item.Product.ID]).Error
		if err != nil {
			err = entity.NewError(err.Error(), http.StatusInternalServerError)
			tx.Rollback()
			return
		}
	}

	err = tx.Commit().Error
	return
}

func (r *repo) pluckProductIDFromCheckoutItems(items []*entity.CheckoutItem) []int64 {
	var result []int64
	for _, item := range items {
		result = append(result, item.Product.ID)
	}
	return result
}

// lock and get product quantity
// return map[int64] where int64 = product id
func (r *repo) lockAndMapProductQuantity(productIDs []int64, tx *gorm.DB) (map[int64]*entity.ProductQuantity, error) {
	var productQuantity []*entity.ProductQuantity
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id in (?)", productIDs).
		Find(&productQuantity).
		Error
	if err != nil {
		return nil, err
	}

	// maping
	result := map[int64]*entity.ProductQuantity{}
	for _, p := range productQuantity {
		result[p.ProductID] = p
	}
	return result, nil
}
