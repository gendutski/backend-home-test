package entity

type MapProductSerialQuantity map[string]int

func (e MapProductSerialQuantity) PluckSerial() []string {
	var result []string
	for serial := range e {
		result = append(result, serial)
	}
	return result
}

type CheckoutItem struct {
	Product       *Product
	Quantity      int
	SubTotalPrice float64
}

type Checkout struct {
	Items      []*CheckoutItem
	TotalItem  int
	TotalPrice float64
}
