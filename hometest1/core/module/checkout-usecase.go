package module

import (
	"net/http"

	"hometest1/core/entity"
	"hometest1/core/repository"
)

type CheckoutUsecase interface {
	Submit(payload entity.MapProductSerialQuantity) (*entity.Checkout, error)
}

type checkoutUsecase struct {
	productRepo repository.ProductRepo
	promoRepo   repository.PromotionRepo
}

func NewCheckoutUsecase(productRepo repository.ProductRepo, promoRepo repository.PromotionRepo) CheckoutUsecase {
	return &checkoutUsecase{productRepo, promoRepo}
}

func (uc *checkoutUsecase) Submit(payload entity.MapProductSerialQuantity) (*entity.Checkout, error) {
	// get products
	products, err := uc.productRepo.GetProductBySerials(payload.PluckSerial())
	if err != nil {
		return nil, entity.NewError(err.Error(), http.StatusInternalServerError)
	}
	if len(products) == 0 {
		return nil, entity.NewError(entity.ProductNotFound, http.StatusBadRequest)
	}

	// get promotions
	promotionMaps, err := uc.promoRepo.GetPromotionByProducts(products)
	if err != nil {
		return nil, entity.NewError(err.Error(), http.StatusInternalServerError)
	}

	// render checkout
	checkout, err := uc.generateCheckout(payload, products, promotionMaps)
	if err != nil {
		return nil, err
	}

	// submit checkout to database
	err = uc.productRepo.SubmitCheckout(checkout)
	if err != nil {
		// repository must handle error with entity.Err
		return nil, err
	}

	return checkout, nil
}

func (uc *checkoutUsecase) generateCheckout(mapQuantity entity.MapProductSerialQuantity, products []*entity.Product, promotionMaps map[int64][]*entity.Promotion) (*entity.Checkout, error) {
	// if product item is free by promo
	// map[int64] = product id, int = number available free items
	var freeProductItem map[int64]int

	var result entity.Checkout

	// loop products
	for _, product := range products {
		var checkoutItem entity.CheckoutItem

		// set quantity
		qty := mapQuantity[product.Serial]

		// init checkout
		checkoutItem.Product = product
		checkoutItem.Quantity = qty
		checkoutItem.SubTotalPrice = float64(qty) * product.Price

		for productID, promos := range promotionMaps {
			if productID == product.ID {
				for _, promo := range promos {
					// the repository should sort promotion types in ascending order
					switch promo.Type {
					case entity.BonusItem:
						freeProductItem = uc.handleBonusItemPromotion(qty, promo, freeProductItem)
					case entity.BuyItemsForReducePrice:
						checkoutItem.SubTotalPrice = uc.handleReducePricePromotion(qty, product, promo)
					case entity.DiscountInPercent:
						checkoutItem.SubTotalPrice = uc.handleDiscountPromotion(qty, checkoutItem.SubTotalPrice, promo)
					default:
						continue
					}
				}
			}
		}

		// set result
		result.Items = append(result.Items, &checkoutItem)
		result.TotalItem += qty
		result.TotalPrice += checkoutItem.SubTotalPrice
	}

	err := uc.handleCheckoutFreeItems(&result, freeProductItem)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// This function calculates the free items that will be obtained
func (uc *checkoutUsecase) handleBonusItemPromotion(quantity int, promo *entity.Promotion, freeProductItem map[int64]int) map[int64]int {
	// if promo product id empty or no match quantity, no free item for this promo
	if promo.PromoProductID == 0 || promo.MatchQuantity == 0 || quantity < promo.MatchQuantity {
		return freeProductItem
	}

	if freeProductItem == nil {
		freeProductItem = make(map[int64]int)
	}

	// number of free item will user get
	numOfFreeItems := (quantity / promo.MatchQuantity) * promo.PromoValue
	freeProductItem[promo.PromoProductID] += numOfFreeItems
	return freeProductItem
}

// This function calculates price reductions that apply multiples
func (uc *checkoutUsecase) handleReducePricePromotion(quantity int, product *entity.Product, promo *entity.Promotion) float64 {
	// if match quantity empty, return original price
	if promo.MatchQuantity <= 0 || quantity < promo.MatchQuantity {
		return product.Price * float64(quantity)
	}

	// get item reduction
	newQuantity := (quantity / promo.MatchQuantity * promo.PromoValue) + (quantity % promo.MatchQuantity)
	return product.Price * float64(newQuantity)
}

// This function calculates the discount price
func (uc *checkoutUsecase) handleDiscountPromotion(quantity int, currentSubTotal float64, promo *entity.Promotion) float64 {
	// promo value for discount in percent value, only process valid value
	if promo.PromoValue < 0 || promo.PromoValue > 100 || quantity < promo.MatchQuantity {
		return currentSubTotal
	}

	return currentSubTotal - (currentSubTotal * float64(promo.PromoValue) / float64(100))
}

// This will handle free items obtained through promotions
// If the item is there, the fee will be deducted, if it is not there it will be added to checkout
func (uc *checkoutUsecase) handleCheckoutFreeItems(checkout *entity.Checkout, freeProductItem map[int64]int) error {
	// check the item in the existing checkout items list
	for _, item := range checkout.Items {
		if freeProductItem[item.Product.ID] > 0 {
			// If the number of items is less than it should be
			if item.Quantity < freeProductItem[item.Product.ID] {
				// reduce checkout total price for current quantity
				checkout.TotalPrice -= float64(item.Quantity) * item.Product.Price
				// add remaining quantity amount
				checkout.TotalItem += freeProductItem[item.Product.ID] - item.Quantity
				item.Quantity = freeProductItem[item.Product.ID]
				item.SubTotalPrice = 0
			} else {
				// if the free items exceed the total items, only reduce the price of the available free items
				priceReduction := float64(freeProductItem[item.Product.ID]) * item.Product.Price
				item.SubTotalPrice = item.SubTotalPrice - priceReduction
				// reduce the sub total price
				checkout.TotalPrice -= priceReduction
			}

			// empty free product item
			freeProductItem[item.Product.ID] = 0
		}
	}

	// if freeProductItem still have quantity, add to checkout
	var productIDs []int64
	for id, qty := range freeProductItem {
		if qty == 0 {
			continue
		}
		productIDs = append(productIDs, id)
	}

	if len(productIDs) == 0 {
		return nil
	}

	// get free product
	products, err := uc.productRepo.GetProductByIDs(productIDs)
	if err != nil {
		return entity.NewError(err.Error(), http.StatusInternalServerError)
	}

	// append to checkout
	for _, product := range products {
		// append new items
		checkout.Items = append(checkout.Items, &entity.CheckoutItem{
			Product:  product,
			Quantity: freeProductItem[product.ID],
		})
		// append checkout total item
		checkout.TotalItem += freeProductItem[product.ID]
		// empty free product item
		freeProductItem[product.ID] = 0
	}

	return nil
}
