package repository_test

import (
	"fmt"
	"hometest4/core/entity"
	"hometest4/core/repository"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := make(map[string]*entity.User)
		repo := repository.NewRepository(user, nil, nil)

		var wg sync.WaitGroup
		var errCh = make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				_, err := repo.CreateUser(fmt.Sprintf("user-%04d", i))
				errCh <- err
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}
		assert.Equal(t, 1000, len(user))
		assert.Empty(t, errs)
	})

	t.Run("duplicate test", func(t *testing.T) {
		user := map[string]*entity.User{
			"odd":  {Username: "odd"},
			"even": {Username: "even"},
		}
		repo := repository.NewRepository(user, nil, nil)

		var wg sync.WaitGroup
		var errCh = make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				username := "even"
				if i%2 == 0 {
					username = "odd"
				}
				_, err := repo.CreateUser(username)
				errCh <- err
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}
		assert.Equal(t, 1000, len(errs))
		assert.Equal(t, 2, len(user))
	})

	t.Run("empty username test", func(t *testing.T) {
		user := make(map[string]*entity.User)
		repo := repository.NewRepository(user, nil, nil)

		var wg sync.WaitGroup
		var errCh = make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- error) {
				defer wg.Done()
				_, err := repo.CreateUser("")
				errCh <- err
			}(i, errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}
		assert.Equal(t, 1000, len(errs))
		assert.Empty(t, user)
	})
}

func Test_ReadUserBalance(t *testing.T) {
	t.Run("register & unregister user test", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman", Balance: 100000},
			"wati":   {Username: "wati", Balance: 200000},
			"kansa":  {Username: "kansa", Balance: 300000},
			"fira":   {Username: "fira", Balance: 400000},
		}
		repo := repository.NewRepository(users, nil, nil)

		var wg sync.WaitGroup
		var resultCh = make(chan map[string]map[string]interface{})
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int, ch chan<- map[string]map[string]interface{}) {
				defer wg.Done()
				result := map[string]map[string]interface{}{}
				for _, user := range []string{"firman", "wati", "kansa", "fira", "umar"} {
					balance, err := repo.ReadUserBalance(user)
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
			for key, usr := range users {
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
		users := map[string]*entity.User{
			"firman": {Username: "firman"},
		}
		topup := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, topup, nil)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TopupUserBalance("firman", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}

		assert.Empty(t, errs)
		assert.Equal(t, float64(10000), users["firman"].Balance)
		assert.Equal(t, 1000, len(topup["firman"]))
	})

	t.Run("unregister user", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman"},
		}
		topup := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, topup, nil)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TopupUserBalance("gendutski", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}

		assert.Equal(t, 1000, len(errs))
		assert.Equal(t, float64(0), users["firman"].Balance)
		assert.Equal(t, 0, len(topup["firman"]))

	})
	t.Run("empty username", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman"},
		}
		topup := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, topup, nil)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TopupUserBalance("", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}

		assert.Equal(t, 1000, len(errs))
		assert.Equal(t, float64(0), users["firman"].Balance)
		assert.Equal(t, 0, len(topup["firman"]))
	})
}

func Test_TransferToUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman", Balance: 10000},
			"wati":   {Username: "wati"},
		}
		transfer := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, nil, transfer)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TransferToUser("firman", "wati", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.Nil(t, err)
		}
		assert.Equal(t, float64(10000), users["wati"].Balance)
		assert.Equal(t, float64(0), users["firman"].Balance)
		for i := 0; i < 1000; i++ {
			assert.Equal(t, []interface{}{
				"wati",
				entity.TrxFlowDebit,
				float64(10),
			}, []interface{}{
				transfer["firman"][i].Username,
				transfer["firman"][i].Flow,
				transfer["firman"][i].Amount,
			})

			assert.Equal(t, []interface{}{
				"firman",
				entity.TrxFlowCredit,
				float64(10),
			}, []interface{}{
				transfer["wati"][i].Username,
				transfer["wati"][i].Flow,
				transfer["wati"][i].Amount,
			})
		}
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman", Balance: 1000},
			"wati":   {Username: "wati"},
		}
		transfer := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, nil, transfer)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TransferToUser("firman", "wati", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		var errs []error
		for err := range errCh {
			if err != nil {
				errs = append(errs, err)
			}
		}
		assert.Equal(t, 900, len(errs))
		assert.Equal(t, float64(1000), users["wati"].Balance)
		assert.Equal(t, float64(0), users["firman"].Balance)
		assert.Equal(t, 100, len(transfer["wati"]))
		assert.Equal(t, 100, len(transfer["firman"]))
		for i := 0; i < 100; i++ {
			assert.Equal(t, []interface{}{
				"wati",
				entity.TrxFlowDebit,
				float64(10),
			}, []interface{}{
				transfer["firman"][i].Username,
				transfer["firman"][i].Flow,
				transfer["firman"][i].Amount,
			})

			assert.Equal(t, []interface{}{
				"firman",
				entity.TrxFlowCredit,
				float64(10),
			}, []interface{}{
				transfer["wati"][i].Username,
				transfer["wati"][i].Flow,
				transfer["wati"][i].Amount,
			})
		}
	})

	t.Run("Destination user not found", func(t *testing.T) {
		users := map[string]*entity.User{
			"firman": {Username: "firman", Balance: 10000},
		}
		transfer := make(map[string][]*entity.Transaction)
		repo := repository.NewRepository(users, nil, transfer)

		var wg sync.WaitGroup
		errCh := make(chan error)
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(ch chan<- error) {
				defer wg.Done()
				errCh <- repo.TransferToUser("firman", "wati", 10)
			}(errCh)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			assert.NotNil(t, err)
		}
		assert.Equal(t, float64(10000), users["firman"].Balance)
		assert.Empty(t, transfer["firman"])
	})
}

func Test_TopTransactionByAmount(t *testing.T) {
	t.Run("have transaction", func(t *testing.T) {
		users := map[string]*entity.User{"firman": {Username: "firman"}}
		transfer := make(map[string][]*entity.Transaction)

		//set transfer data
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("user-%04d", i)
			flow := entity.TrxFlowCredit
			if i%2 == 0 {
				flow = entity.TrxFlowDebit
			}
			transfer["firman"] = append(transfer["firman"], &entity.Transaction{
				Username: key,
				Amount:   float64((i + 1) * 10),
				Flow:     flow,
			})
		}

		repo := repository.NewRepository(users, nil, transfer)
		result, err := repo.TopTransactionByAmount("firman")
		assert.Nil(t, err)
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
	})

	t.Run("no transaction", func(t *testing.T) {
		users := map[string]*entity.User{"firman": {Username: "firman"}}
		transfer := make(map[string][]*entity.Transaction)

		repo := repository.NewRepository(users, nil, transfer)
		result, err := repo.TopTransactionByAmount("firman")
		assert.Nil(t, err)
		assert.Empty(t, result)
	})
}

func Test_TopTransactionBySumAmount(t *testing.T) {
	t.Run("have transaction", func(t *testing.T) {
		users := map[string]*entity.User{"firman": {Username: "firman"}}
		transfer := make(map[string][]*entity.Transaction)

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
			transfer["firman"] = append(transfer["firman"], &entity.Transaction{
				Username: key,
				Amount:   amount,
				Flow:     flow,
			})
		}

		repo := repository.NewRepository(users, nil, transfer)
		result, err := repo.TopTransactionBySumAmount("firman")
		assert.Nil(t, err)
		assert.Equal(t, []*entity.OverallTopTransactionResponse{
			{Username: "user-02", TransactionValue: validAmount},
		}, result)
	})

	t.Run("no transaction", func(t *testing.T) {
		users := map[string]*entity.User{"firman": {Username: "firman"}}
		transfer := make(map[string][]*entity.Transaction)

		repo := repository.NewRepository(users, nil, transfer)
		result, err := repo.TopTransactionBySumAmount("firman")
		assert.Nil(t, err)
		assert.Empty(t, result)
	})
}
