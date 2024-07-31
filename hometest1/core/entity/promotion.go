package entity

import (
	"time"

	"gorm.io/gorm"
)

type PromotionType int

const (
	UndefinedType PromotionType = iota
	BonusItem
	BuyItemsForReducePrice
	DiscountInPercent
	FreeItem
)

type Promotion struct {
	ID             int64
	Type           PromotionType
	ProductID      int64
	MatchQuantity  int
	PromoValue     int
	PromoProductID int64
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
