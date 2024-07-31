package repository

import (
	"hometest4/config"
	"hometest4/core/entity"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	GetUserDetail(username string) (*entity.User, error)
	CreateUser(username string) (*entity.User, error)
	ReadUserBalance(username string) (float64, error)
	TopupUserBalance(username string, amount float64) error
	TransferToUser(username string, toUsername string, amount float64) error
	TopTransactionByAmount(username string) ([]*entity.TopTransactionResponse, error)
	TopTransactionBySumAmount(username string) ([]*entity.OverallTopTransactionResponse, error)
}

func NewRepository(
	user map[string]*entity.User,
	topup map[string][]*entity.Transaction,
	transfer map[string][]*entity.Transaction,
) Repository {
	return &repo{
		user:     user,
		topup:    topup,
		transfer: transfer,
	}
}

type repo struct {
	rw       sync.RWMutex
	user     map[string]*entity.User
	topup    map[string][]*entity.Transaction
	transfer map[string][]*entity.Transaction
}

func (r *repo) GetUserDetail(username string) (*entity.User, error) {
	if username == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	r.rw.RLock()
	defer r.rw.RUnlock()
	user, ok := r.user[username]
	if !ok {
		return nil, echo.NewHTTPError(http.StatusNotFound)
	}
	return user, nil
}

func (r *repo) CreateUser(username string) (*entity.User, error) {
	// get user
	_, err := r.GetUserDetail(username)
	if err == nil {
		return nil, echo.NewHTTPError(http.StatusConflict, "Username already exists")
	}
	if herr, ok := err.(*echo.HTTPError); !ok || herr.Code != http.StatusNotFound {
		return nil, err
	}

	// set user
	user := &entity.User{
		Username: username,
	}
	r.rw.Lock()
	r.user[username] = user
	r.rw.Unlock()

	return user, nil
}

func (r *repo) ReadUserBalance(username string) (float64, error) {
	if username == "" {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	r.rw.RLock()
	defer r.rw.RUnlock()
	user, ok := r.user[username]
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized)
	}
	return user.Balance, nil
}

func (r *repo) TopupUserBalance(username string, amount float64) error {
	if username == "" || amount <= 0 {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	r.rw.Lock()
	defer r.rw.Unlock()

	// unregistered user
	if _, ok := r.user[username]; !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	// set user balance
	r.user[username].Balance += amount
	r.topup[username] = append(r.topup[username], &entity.Transaction{
		Username: username,
		Amount:   amount,
		Flow:     entity.TrxFlowDebit,
		TrxDate:  time.Now(),
	})
	return nil
}

func (r *repo) TransferToUser(username string, toUsername string, amount float64) error {
	if username == "" || toUsername == "" || amount <= 0 {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// check current user
	user, err := r.GetUserDetail(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	// check target user
	_, err = r.GetUserDetail(toUsername)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Destination user not found")
	}

	r.rw.Lock()
	defer r.rw.Unlock()
	// check user balance
	if user.Balance-amount < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Insufficient balance")
	}
	// set user balance
	r.user[username].Balance -= amount
	// set user transaction data
	r.transfer[username] = append(r.transfer[username], &entity.Transaction{
		Username: toUsername,
		Amount:   amount,
		Flow:     entity.TrxFlowDebit,
		TrxDate:  time.Now(),
	})

	// set target user balance
	r.user[toUsername].Balance += amount
	// set target user transaction data
	r.transfer[toUsername] = append(r.transfer[toUsername], &entity.Transaction{
		Username: username,
		Amount:   amount,
		Flow:     entity.TrxFlowCredit,
		TrxDate:  time.Now(),
	})
	return nil
}

func (r *repo) TopTransactionByAmount(username string) ([]*entity.TopTransactionResponse, error) {
	// validate user
	if _, err := r.GetUserDetail(username); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	// set result
	var result []*entity.TopTransactionResponse
	for _, trx := range r.transfer[username] {
		amount := trx.Amount
		if trx.Flow == entity.TrxFlowDebit {
			amount *= -1
		}

		result = append(result, &entity.TopTransactionResponse{
			Username:         trx.Username,
			TransactionValue: amount,
		})
	}

	// sort result
	sort.Slice(result, func(i, j int) bool {
		return math.Abs(result[i].TransactionValue) > math.Abs(result[j].TransactionValue)
	})
	if len(result) > config.TopTransactionLimit {
		return result[:config.TopTransactionLimit], nil
	}
	return result, nil
}

func (r *repo) TopTransactionBySumAmount(username string) ([]*entity.OverallTopTransactionResponse, error) {
	// validate user
	if _, err := r.GetUserDetail(username); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError)
	}

	r.rw.RLock()
	defer r.rw.RUnlock()

	// sum per user
	maper := make(map[string]float64)
	for _, trx := range r.transfer[username] {
		if trx.Flow == entity.TrxFlowDebit {
			maper[trx.Username] += trx.Amount
		}
	}

	// set result
	var result []*entity.OverallTopTransactionResponse
	for usr, amount := range maper {
		result = append(result, &entity.OverallTopTransactionResponse{
			Username:         usr,
			TransactionValue: amount,
		})
	}

	// sort result
	sort.Slice(result, func(i, j int) bool {
		return result[i].TransactionValue > result[j].TransactionValue
	})
	if len(result) > config.TopTransactionLimit {
		return result[:config.TopTransactionLimit], nil
	}
	return result, nil
}
