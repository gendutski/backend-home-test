package service_test

import (
	"errors"
	"fmt"
	"hometest4/config"
	"hometest4/core/entity"
	"hometest4/core/repository"
	"hometest4/core/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type data struct {
	User     map[string]*entity.User
	Topup    map[string][]*entity.Transaction
	Transfer map[string][]*entity.Transaction
}

func initService(theData *data) service.Service {
	repo := repository.NewRepository(theData.User, theData.Topup, theData.Transfer)
	service := service.NewService(repo)
	return service
}

func validateTokenString(tokenStr string, validUsername string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetJwtSecret()), nil
	})
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if holder, ok := claims[config.TokenUsernameField]; ok {
			if username, ok := holder.(string); ok && username == validUsername {
				return true
			}
		}
	}
	return false
}

func Test_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		theData := &data{
			User: make(map[string]*entity.User),
		}
		service := initService(theData)

		var wg sync.WaitGroup
		var respCh = make(chan map[string]interface{})
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- map[string]interface{}) {
				defer wg.Done()
				username := fmt.Sprintf("user-%04d", i)
				token, err := service.CreateUser(&entity.RegisterPayload{
					Username: username,
				})
				ch <- map[string]interface{}{
					"Username": username,
					"Token":    token,
					"Error":    err,
				}
			}(i, respCh)
		}

		go func() {
			wg.Wait()
			close(respCh)
		}()

		var tokens []string
		var err []error
		for resp := range respCh {
			if resp["Error"] != nil {
				err = append(err, resp["Error"].(error))
			} else {
				if validateTokenString(resp["Token"].(string), resp["Username"].(string)) {
					tokens = append(tokens, resp["Token"].(string))
				} else {
					err = append(err, errors.New("invalid token"))
				}
			}
		}

		assert.Empty(t, err)
		assert.Equal(t, 1000, len(tokens))
		assert.Equal(t, 1000, len(theData.User))
	})

	t.Run("duplicate username", func(t *testing.T) {
		theData := &data{
			User: map[string]*entity.User{
				"even": {Username: "even"},
				"odd":  {Username: "odd"},
			},
		}
		service := initService(theData)

		var wg sync.WaitGroup
		var respCh = make(chan map[string]interface{})
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- map[string]interface{}) {
				defer wg.Done()
				username := "even"
				if i%2 == 0 {
					username = "odd"
				}
				token, err := service.CreateUser(&entity.RegisterPayload{
					Username: username,
				})
				ch <- map[string]interface{}{
					"Username": username,
					"Token":    token,
					"Error":    err,
				}
			}(i, respCh)
		}

		go func() {
			wg.Wait()
			close(respCh)
		}()

		var tokens []string
		var err []error
		for resp := range respCh {
			if resp["Error"] != nil {
				herr, ok := resp["Error"].(*echo.HTTPError)
				if !ok {
					continue
				}
				// only 409 error accepted
				if herr.Code == http.StatusConflict {
					err = append(err, resp["Error"].(error))
				}
			}
		}

		assert.Equal(t, 1000, len(err))
		assert.Empty(t, tokens)
		assert.Equal(t, 2, len(theData.User))
	})

	t.Run("bad request", func(t *testing.T) {
		theData := &data{}
		service := initService(theData)
		var wg sync.WaitGroup
		var respCh = make(chan map[string]interface{})
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- map[string]interface{}) {
				defer wg.Done()
				var payload entity.RegisterPayload
				if i%2 == 0 {
					payload.Username = "inval!d username"
				} else if i%3 == 0 {
					payload.Username = "long-username-exceeded-maximum-30-characters-for-userame"
				}

				token, err := service.CreateUser(&payload)
				ch <- map[string]interface{}{
					"Token": token,
					"Error": err,
				}
			}(i, respCh)
		}

		go func() {
			wg.Wait()
			close(respCh)
		}()

		var tokens []string
		var err []error
		for resp := range respCh {
			if resp["Error"] != nil {
				herr, ok := resp["Error"].(*echo.HTTPError)
				if !ok {
					continue
				}
				// only 400 error accepted
				if herr.Code == http.StatusBadRequest {
					err = append(err, resp["Error"].(error))
				}
			}
		}

		assert.Equal(t, 1000, len(err))
		assert.Empty(t, tokens)
		assert.Empty(t, theData.User)
	})
}

func Test_ReadBalance(t *testing.T) {
	theData := &data{
		User: map[string]*entity.User{
			"firman": {Username: "firman", Balance: 100000},
			"wati":   {Username: "wati", Balance: 200000},
			"kansa":  {Username: "kansa", Balance: 300000},
			"fira":   {Username: "fira", Balance: 400000},
		},
	}
	service := initService(theData)
	e := echo.New()

	t.Run("register & unregister user test", func(t *testing.T) {
		var wg sync.WaitGroup
		var resultCh = make(chan map[string]map[string]interface{})
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- map[string]map[string]interface{}) {
				defer wg.Done()
				result := map[string]map[string]interface{}{}
				for _, user := range []string{"firman", "wati", "kansa", "fira", "umar"} {
					req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
					rec := httptest.NewRecorder()
					eContext := e.NewContext(req, rec)
					eContext.Set(config.UserContextKey, user)

					balance, err := service.ReadBalance(eContext)
					result[user] = map[string]interface{}{
						"balance": balance,
						"error":   err,
					}
				}
				resultCh <- result
			}(i, resultCh)
		}

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for result := range resultCh {
			// register users
			for key, usr := range theData.User {
				assert.Equal(t, usr.Balance, result[key]["balance"])
				assert.Nil(t, result[key]["error"])
			}
			// unregister users
			assert.Empty(t, result["umar"]["balance"])
			assert.NotNil(t, result["umar"]["error"])
		}
	})
}

func Test_TopupUserBalance(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 0},
			},
			Topup: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				errCh <- service.TopupBalance(eContext, &entity.BalanceTopupPayload{
					Amount: 10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.Nil(t, err)
		}
		assert.Equal(t, float64(10000), theData.User["firman"].Balance)
	})

	t.Run("Invalid topup amount", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 0},
			},
			Topup: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				var amount float64 = -100
				if i%2 == 0 {
					amount = 10000000
				} else if i%3 == 0 {
					amount = 0
				}
				errCh <- service.TopupBalance(eContext, &entity.BalanceTopupPayload{
					Amount: amount,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, herr.Code)
		}
		assert.Equal(t, float64(0), theData.User["firman"].Balance)
	})

	t.Run("unauthorized", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 0},
			},
			Topup: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)

				errCh <- service.TopupBalance(eContext, &entity.BalanceTopupPayload{
					Amount: 10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, herr.Code)
		}
		assert.Equal(t, float64(0), theData.User["firman"].Balance)
	})
}

func Test_Transfer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 10000},
				"wati":   {Username: "wati", Balance: 0},
			},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				errCh <- service.Transfer(eContext, &entity.TransferPayload{
					ToUsername: "wati",
					Amount:     10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.Nil(t, err)
		}
		assert.Equal(t, float64(0), theData.User["firman"].Balance)
		assert.Equal(t, float64(10000), theData.User["wati"].Balance)
		for i := 0; i < 1000; i++ {
			assert.Equal(t, []interface{}{
				"wati",
				entity.TrxFlowDebit,
				float64(10),
			}, []interface{}{
				theData.Transfer["firman"][i].Username,
				theData.Transfer["firman"][i].Flow,
				theData.Transfer["firman"][i].Amount,
			})

			assert.Equal(t, []interface{}{
				"firman",
				entity.TrxFlowCredit,
				float64(10),
			}, []interface{}{
				theData.Transfer["wati"][i].Username,
				theData.Transfer["wati"][i].Flow,
				theData.Transfer["wati"][i].Amount,
			})
		}
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 1000},
				"wati":   {Username: "wati", Balance: 0},
			},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				errCh <- service.Transfer(eContext, &entity.TransferPayload{
					ToUsername: "wati",
					Amount:     10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				if herr, ok := err.(*echo.HTTPError); ok && herr.Code == http.StatusBadRequest {
					errs = append(errs, err)
				}
			}
		}
		assert.Equal(t, 900, len(errs))
		assert.Equal(t, float64(0), theData.User["firman"].Balance)
		assert.Equal(t, float64(1000), theData.User["wati"].Balance)
		assert.Equal(t, 100, len(theData.Transfer["wati"]))
		assert.Equal(t, 100, len(theData.Transfer["firman"]))
		for i := 0; i < 100; i++ {
			assert.Equal(t, []interface{}{
				"wati",
				entity.TrxFlowDebit,
				float64(10),
			}, []interface{}{
				theData.Transfer["firman"][i].Username,
				theData.Transfer["firman"][i].Flow,
				theData.Transfer["firman"][i].Amount,
			})

			assert.Equal(t, []interface{}{
				"firman",
				entity.TrxFlowCredit,
				float64(10),
			}, []interface{}{
				theData.Transfer["wati"][i].Username,
				theData.Transfer["wati"][i].Flow,
				theData.Transfer["wati"][i].Amount,
			})
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 10000},
				"wati":   {Username: "wati", Balance: 0},
			},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)

				errCh <- service.Transfer(eContext, &entity.TransferPayload{
					ToUsername: "wati",
					Amount:     10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, herr.Code)
		}
	})

	t.Run("Destination user not found", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User: map[string]*entity.User{
				"firman": {Username: "firman", Balance: 10000},
			},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				errCh <- service.Transfer(eContext, &entity.TransferPayload{
					ToUsername: "wati",
					Amount:     10,
				})
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusNotFound, herr.Code)
		}
		assert.Equal(t, float64(10000), theData.User["firman"].Balance)
		assert.Empty(t, theData.Transfer["firman"])
	})
}

func Test_TopTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		var resultCh = make(chan []*entity.TopTransactionResponse)

		theData := &data{
			User:     map[string]*entity.User{"firman": {Username: "firman"}},
			Transfer: make(map[string][]*entity.Transaction),
		}
		//set transfer data
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("user-%04d", i)
			flow := entity.TrxFlowCredit
			if i%2 == 0 {
				flow = entity.TrxFlowDebit
			}
			theData.Transfer["firman"] = append(theData.Transfer["firman"], &entity.Transaction{
				Username: key,
				Amount:   float64((i + 1) * 10),
				Flow:     flow,
			})
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- []*entity.TopTransactionResponse) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				result, _ := service.TopTransaction(eContext)
				resultCh <- result
			}(i, resultCh)
		}

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for result := range resultCh {
			assert.Equal(t, 10, len(result))
			assert.Equal(t, []*entity.TopTransactionResponse{
				{Username: "user-0099", TransactionValue: 1000},
				{Username: "user-0098", TransactionValue: -990},
				{Username: "user-0097", TransactionValue: 980},
				{Username: "user-0096", TransactionValue: -970},
				{Username: "user-0095", TransactionValue: 960},
				{Username: "user-0094", TransactionValue: -950},
				{Username: "user-0093", TransactionValue: 940},
				{Username: "user-0092", TransactionValue: -930},
				{Username: "user-0091", TransactionValue: 920},
				{Username: "user-0090", TransactionValue: -910},
			}, result)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User:     map[string]*entity.User{"firman": {Username: "firman"}},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)

				_, err := service.TopTransaction(eContext)
				errCh <- err
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, herr.Code)
		}
	})
}

func Test_OverallTopTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wg sync.WaitGroup
		var resultCh = make(chan []*entity.OverallTopTransactionResponse)

		theData := &data{
			User:     map[string]*entity.User{"firman": {Username: "firman"}},
			Transfer: make(map[string][]*entity.Transaction),
		}
		//set transfer data
		var validAmount float64
		for i := 0; i < 100; i++ {
			key := "user-01"
			flow := entity.TrxFlowCredit
			amount := float64((i + 1) * 20)
			if i%2 == 0 {
				key = "user-02"
				amount = float64((i + 1) * 10)
				flow = entity.TrxFlowDebit
				validAmount += amount
			}
			theData.Transfer["firman"] = append(theData.Transfer["firman"], &entity.Transaction{
				Username: key,
				Amount:   amount,
				Flow:     flow,
			})
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- []*entity.OverallTopTransactionResponse) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)
				eContext.Set(config.UserContextKey, "firman")

				result, _ := service.OverallTopTransaction(eContext)
				resultCh <- result
			}(i, resultCh)
		}

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for result := range resultCh {
			assert.Equal(t, 1, len(result))
			assert.Equal(t, []*entity.OverallTopTransactionResponse{
				{Username: "user-02", TransactionValue: validAmount},
			}, result)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		var wg sync.WaitGroup
		var errCh = make(chan error)

		theData := &data{
			User:     map[string]*entity.User{"firman": {Username: "firman"}},
			Transfer: make(map[string][]*entity.Transaction),
		}
		service := initService(theData)
		e := echo.New()

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				rec := httptest.NewRecorder()
				eContext := e.NewContext(req, rec)

				_, err := service.TopTransaction(eContext)
				errCh <- err
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
			herr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, herr.Code)
		}
	})
}
