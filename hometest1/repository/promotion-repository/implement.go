package promotionrepository

import (
	"hometest1/core/entity"
	"hometest1/core/repository"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.PromotionRepo {
	return &repo{db}
}

func (r *repo) GetPromotionByProducts(products []*entity.Product) (map[int64][]*entity.Promotion, error) {
	// pluck product id
	var ids []int64
	for _, p := range products {
		ids = append(ids, p.ID)
	}

	// get promotions by product id
	var promotions []*entity.Promotion
	err := r.db.Where("product_id in (?)", ids).Order("product_id asc, type asc").Find(&promotions).Error
	if err != nil {
		return nil, err
	}

	// maping product
	result := map[int64][]*entity.Promotion{}
	for _, promo := range promotions {
		result[promo.ProductID] = append(result[promo.ProductID], promo)
	}

	return result, nil
}
