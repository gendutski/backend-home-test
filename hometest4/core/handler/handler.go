package handler

import (
	"hometest4/core/entity"
	"hometest4/core/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHandler(srv service.Service) *Handler {
	return &Handler{srv}
}

type Handler struct {
	srv service.Service
}

func (h *Handler) BalanceRead(e echo.Context) error {
	balance, err := h.srv.ReadBalance(e)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, entity.BalanceReadResponse{
		Balance: balance,
	})
}

func (h *Handler) Transfer(e echo.Context) error {
	payload := new(entity.TransferPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	err = h.srv.Transfer(e, payload)
	if err != nil {
		return err
	}
	return e.String(http.StatusNoContent, "Transfer success")
}

func (h *Handler) ListOverallTopTrxUsersByValue(e echo.Context) error {
	result, err := h.srv.OverallTopTransaction(e)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, result)
}

func (h *Handler) RegisterUser(e echo.Context) error {
	// bind payload
	payload := new(entity.RegisterPayload)
	err := e.Bind(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// register user
	token, err := h.srv.CreateUser(payload)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, entity.RegisterResponse{
		Token: token,
	})
}

func (h *Handler) BalanceTopup(e echo.Context) error {
	payload := new(entity.BalanceTopupPayload)
	err := e.Bind(payload)
	if err != nil {
		return err
	}
	err = h.srv.TopupBalance(e, payload)
	if err != nil {
		return err
	}
	return e.String(http.StatusNoContent, "Topup successful")
}

func (h *Handler) TopTrxForUser(e echo.Context) error {
	result, err := h.srv.TopTransaction(e)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, result)
}
