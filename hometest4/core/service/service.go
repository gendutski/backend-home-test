package service

import (
	"hometest4/config"
	"hometest4/core/entity"
	"hometest4/core/repository"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Service interface {
	CreateUser(payload *entity.RegisterPayload) (string, error)
	ReadBalance(e echo.Context) (float64, error)
	TopupBalance(e echo.Context, payload *entity.BalanceTopupPayload) error
	Transfer(e echo.Context, payload *entity.TransferPayload) error
	TopTransaction(e echo.Context) ([]*entity.TopTransactionResponse, error)
	OverallTopTransaction(e echo.Context) ([]*entity.OverallTopTransactionResponse, error)
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo:      repo,
		validator: config.InitValidator(),
	}
}

type service struct {
	repo      repository.Repository
	validator *config.Validator
}

func (s *service) CreateUser(payload *entity.RegisterPayload) (string, error) {
	// validate payload
	err := s.validator.Validate(payload, "Bad Request")
	if err != nil {
		return "", err
	}

	// get user from data
	user, err := s.repo.CreateUser(payload.Username)
	if err != nil {
		return "", err
	}

	// parse token
	claims := jwt.MapClaims{
		config.TokenUsernameField: user.Username,
		"exp":                     time.Now().Add(config.TokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.GetJwtSecret()))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (s *service) ReadBalance(e echo.Context) (float64, error) {
	username, ok := e.Get(config.UserContextKey).(string)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized)
	}

	return s.repo.ReadUserBalance(username)
}

func (s *service) TopupBalance(e echo.Context, payload *entity.BalanceTopupPayload) error {
	// validate payload
	err := s.validator.Validate(payload, "Invalid topup amount")
	if err != nil {
		return err
	}
	// validate auth
	username, ok := e.Get(config.UserContextKey).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	return s.repo.TopupUserBalance(username, payload.Amount)
}

func (s *service) Transfer(e echo.Context, payload *entity.TransferPayload) error {
	// validate payload
	err := s.validator.Validate(payload, "Insufficient balance")
	if err != nil {
		return err
	}
	// validate auth
	username, ok := e.Get(config.UserContextKey).(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	return s.repo.TransferToUser(username, payload.ToUsername, payload.Amount)
}

func (s *service) TopTransaction(e echo.Context) ([]*entity.TopTransactionResponse, error) {
	// validate auth
	username, ok := e.Get(config.UserContextKey).(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}
	return s.repo.TopTransactionByAmount(username)
}

func (s *service) OverallTopTransaction(e echo.Context) ([]*entity.OverallTopTransactionResponse, error) {
	// validate auth
	username, ok := e.Get(config.UserContextKey).(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}
	return s.repo.TopTransactionBySumAmount(username)
}
