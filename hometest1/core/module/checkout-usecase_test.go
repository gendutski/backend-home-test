package module_test

import (
	"testing"
	"time"

	"hometest1/core/entity"
	"hometest1/core/module"
	repomocks "hometest1/core/repository/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func initCheckoutUC(ctrl *gomock.Controller) (module.CheckoutUsecase, *repomocks.MockProductRepo, *repomocks.MockPromotionRepo) {
	productRepo := repomocks.NewMockProductRepo(ctrl)
	promoRepo := repomocks.NewMockPromotionRepo(ctrl)

	return module.NewCheckoutUsecase(productRepo, promoRepo), productRepo, promoRepo
}

func Test_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, productRepo, promoRepo := initCheckoutUC(ctrl)

	dayCreated, _ := time.Parse("2006-01-02", "2023-05-16")
	products := []*entity.Product{
		{ID: 1, Serial: "120P90", Name: "Google Home", Price: 49.99, UpdatedAt: dayCreated},
		{ID: 2, Serial: "43N23P", Name: "MacBook Pro", Price: 5399.99, UpdatedAt: dayCreated},
		{ID: 3, Serial: "A304SD", Name: "Alexa Speaker", Price: 109.50, UpdatedAt: dayCreated},
		{ID: 4, Serial: "234234", Name: "Raspberry Pi B", Price: 30.00, UpdatedAt: dayCreated},
	}
	promotions := []*entity.Promotion{
		{ID: 1, Type: 1, ProductID: 2, MatchQuantity: 1, PromoValue: 1, PromoProductID: 4, UpdatedAt: dayCreated},
		{ID: 1, Type: 2, ProductID: 1, MatchQuantity: 3, PromoValue: 2, PromoProductID: 0, UpdatedAt: dayCreated},
		{ID: 3, Type: 3, ProductID: 3, MatchQuantity: 3, PromoValue: 10, PromoProductID: 0, UpdatedAt: dayCreated},
	}

	t.Run("Scanned Items: MacBook Pro, Raspberry Pi B", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"43N23P": 1, "234234": 1}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[1], products[3],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[1], products[3],
		}).Return(map[int64][]*entity.Promotion{
			2: {promotions[0]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[1],
					Quantity:      1,
					SubTotalPrice: 5399.99,
				},
				{
					Product:       products[3],
					Quantity:      1,
					SubTotalPrice: 0,
				},
			},
			TotalItem:  2,
			TotalPrice: 5399.99,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: MacBook Pro, 2 Raspberry Pi B", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"43N23P": 1, "234234": 2}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[1], products[3],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[1], products[3],
		}).Return(map[int64][]*entity.Promotion{
			2: {promotions[0]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[1],
					Quantity:      1,
					SubTotalPrice: 5399.99,
				},
				{
					Product:       products[3],
					Quantity:      2,
					SubTotalPrice: 30,
				},
			},
			TotalItem:  3,
			TotalPrice: 5399.99 + 30,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: MacBook Pro, without Raspberry Pi B", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"43N23P": 1}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[1],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[1],
		}).Return(map[int64][]*entity.Promotion{
			2: {promotions[0]},
		}, nil).Times(1)
		productRepo.EXPECT().GetProductByIDs([]int64{4}).Return([]*entity.Product{
			products[3],
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[1],
					Quantity:      1,
					SubTotalPrice: 5399.99,
				},
				{
					Product:       products[3],
					Quantity:      1,
					SubTotalPrice: 0,
				},
			},
			TotalItem:  2,
			TotalPrice: 5399.99,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: 2 MacBook Pro, 1 Raspberry Pi B", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"43N23P": 2, "234234": 1}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[1], products[3],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[1], products[3],
		}).Return(map[int64][]*entity.Promotion{
			2: {promotions[0]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[1],
					Quantity:      2,
					SubTotalPrice: 5399.99 * 2,
				},
				{
					Product:       products[3],
					Quantity:      2,
					SubTotalPrice: 0,
				},
			},
			TotalItem:  4,
			TotalPrice: 5399.99 * 2,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: Google Home, Google Home, Google Home", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"120P90": 3}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[0],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[0],
		}).Return(map[int64][]*entity.Promotion{
			1: {promotions[1]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[0],
					Quantity:      3,
					SubTotalPrice: 49.99 * 2,
				},
			},
			TotalItem:  3,
			TotalPrice: 49.99 * 2,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: 6 Google Home", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"120P90": 6}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[0],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[0],
		}).Return(map[int64][]*entity.Promotion{
			1: {promotions[1]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[0],
					Quantity:      6,
					SubTotalPrice: 49.99 * 4,
				},
			},
			TotalItem:  6,
			TotalPrice: 49.99 * 4,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: 4 Google Home", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"120P90": 4}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[0],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[0],
		}).Return(map[int64][]*entity.Promotion{
			1: {promotions[1]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[0],
					Quantity:      4,
					SubTotalPrice: 49.99 * 3,
				},
			},
			TotalItem:  4,
			TotalPrice: 49.99 * 3,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: Alexa Speaker, Alexa Speaker, Alexa Speaker", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"A304SD": 3}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[2],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[2],
		}).Return(map[int64][]*entity.Promotion{
			3: {promotions[2]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[2],
					Quantity:      3,
					SubTotalPrice: (109.50 * 3) - (109.50 * 3 * 10 / 100),
				},
			},
			TotalItem:  3,
			TotalPrice: (109.50 * 3) - (109.50 * 3 * 10 / 100),
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: 4 Alexa Speaker", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"A304SD": 4}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[2],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[2],
		}).Return(map[int64][]*entity.Promotion{
			3: {promotions[2]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[2],
					Quantity:      4,
					SubTotalPrice: (109.50 * 4) - (109.50 * 4 * 10 / 100),
				},
			},
			TotalItem:  4,
			TotalPrice: (109.50 * 4) - (109.50 * 4 * 10 / 100),
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})

	t.Run("Scanned Items: 2 Alexa Speaker (don't get discount)", func(t *testing.T) {
		payload := entity.MapProductSerialQuantity{"A304SD": 2}
		productRepo.EXPECT().GetProductBySerials(gomock.Any()).Return([]*entity.Product{
			products[2],
		}, nil).Times(1)
		promoRepo.EXPECT().GetPromotionByProducts([]*entity.Product{
			products[2],
		}).Return(map[int64][]*entity.Promotion{
			3: {promotions[2]},
		}, nil).Times(1)

		checkout := &entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:       products[2],
					Quantity:      2,
					SubTotalPrice: 109.50 * 2,
				},
			},
			TotalItem:  2,
			TotalPrice: 109.50 * 2,
		}
		productRepo.EXPECT().SubmitCheckout(checkout).Return(nil).Times(1)

		resp, err := svc.Submit(payload)
		assert.Nil(t, err)
		assert.Equal(t, checkout, resp)
	})
}
