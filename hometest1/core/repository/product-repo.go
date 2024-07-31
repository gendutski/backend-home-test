package repository

import "hometest1/core/entity"

type ProductRepo interface {
	GetProductBySerials(serials []string) ([]*entity.Product, error)
	GetProductByIDs(ids []int64) ([]*entity.Product, error)
	SubmitCheckout(payload *entity.Checkout) error
}
