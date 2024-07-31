package handler

import (
	"net/http"

	"hometest1/core/entity"
	"hometest1/core/module"

	"github.com/labstack/echo/v4"
)

type CheckoutHandler struct {
	checkoutUC module.CheckoutUsecase
}

func NewCheckoutHandler(checkoutUC module.CheckoutUsecase) *CheckoutHandler {
	return &CheckoutHandler{checkoutUC}
}

type payload struct {
	ProductSerials []string `json:"productSerials" validate:"required"`
}

type responseItem struct {
	Serial   string  `json:"serial"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	SubTotal float64 `json:"subTotal"`
}

type response struct {
	Items      []*responseItem `json:"items"`
	TotalItems int             `json:"totalItems"`
	TotalPrice float64         `json:"totalPrice"`
}

func (h *CheckoutHandler) Submit(c echo.Context) error {
	p := new(payload)
	// bind json payload
	if err := c.Bind(p); err != nil {
		return err
	}
	// validate payload
	if err := c.Validate(p); err != nil {
		return err
	}

	// map payload
	mapPayload := make(map[string]int)
	for _, serial := range p.ProductSerials {
		mapPayload[serial]++
	}

	resp, err := h.checkoutUC.Submit(mapPayload)
	if err != nil {
		return err
	}

	return h.parseToResponse(resp, c)
}

func (h *CheckoutHandler) parseToResponse(p *entity.Checkout, c echo.Context) error {
	result := response{
		TotalItems: p.TotalItem,
		TotalPrice: p.TotalPrice,
	}

	for _, item := range p.Items {
		result.Items = append(result.Items, &responseItem{
			Serial:   item.Product.Serial,
			Name:     item.Product.Name,
			Quantity: item.Quantity,
			Price:    item.Product.Price,
			SubTotal: item.SubTotalPrice,
		})
	}

	return c.JSON(http.StatusOK, result)

}
