package productrepository_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"hometest1/core/entity"
	"hometest1/core/repository"
	productrepository "hometest1/repository/product-repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func initRepo(db *sql.DB, mock sqlmock.Sqlmock) (repository.ProductRepo, error) {
	mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.25-log"))
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logger.Info)),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return productrepository.New(gdb), nil
}

func Test_GetProductBySerials(t *testing.T) {
	// mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	defer db.Close()

	// init repo
	repo, err := initRepo(db, mock)
	if err != nil {
		t.Errorf("error initRepo: %s", err.Error())
		return
	}
	dayCreated, _ := time.Parse("2006-01-02", "2023-05-16")

	t.Run("positive", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "serial", "name", "price", "updated_at"}).
			AddRow(1, "120P90", "Google Home", 49.99, dayCreated).
			AddRow(3, "A304SD", "Alexa Speaker", 109.50, dayCreated)

		mock.
			ExpectQuery(regexp.QuoteMeta("SELECT * FROM `product` WHERE serial in (?,?)")).
			WithArgs("120P90", "A304SD").
			WillReturnRows(rows)

		resp, err := repo.GetProductBySerials([]string{"120P90", "A304SD"})
		assert.Nil(t, err)
		assert.Equal(t, []*entity.Product{
			{ID: 1, Serial: "120P90", Name: "Google Home", Price: 49.99, UpdatedAt: dayCreated},
			{ID: 3, Serial: "A304SD", Name: "Alexa Speaker", Price: 109.50, UpdatedAt: dayCreated},
		}, resp)
	})
}

func Test_GetProductByIDs(t *testing.T) {
	// mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	defer db.Close()

	// init repo
	repo, err := initRepo(db, mock)
	if err != nil {
		t.Errorf("error initRepo: %s", err.Error())
		return
	}
	dayCreated, _ := time.Parse("2006-01-02", "2023-05-16")

	t.Run("positive", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "serial", "name", "price", "updated_at"}).
			AddRow(1, "120P90", "Google Home", 49.99, dayCreated).
			AddRow(3, "A304SD", "Alexa Speaker", 109.50, dayCreated)

		mock.
			ExpectQuery(regexp.QuoteMeta("SELECT * FROM `product` WHERE id in (?,?)")).
			WithArgs(1, 3).
			WillReturnRows(rows)

		resp, err := repo.GetProductByIDs([]int64{1, 3})
		assert.Nil(t, err)
		assert.Equal(t, []*entity.Product{
			{ID: 1, Serial: "120P90", Name: "Google Home", Price: 49.99, UpdatedAt: dayCreated},
			{ID: 3, Serial: "A304SD", Name: "Alexa Speaker", Price: 109.50, UpdatedAt: dayCreated},
		}, resp)
	})
}

func Test_SubmitCheckout(t *testing.T) {
	// mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	defer db.Close()

	// init repo
	repo, err := initRepo(db, mock)
	if err != nil {
		t.Errorf("error initRepo: %s", err.Error())
		return
	}
	dayCreated, _ := time.Parse("2006-01-02", "2023-05-16")

	t.Run("positive, item quantity is sufficient", func(t *testing.T) {
		mock.ExpectBegin()

		// lock for update product_quantity
		rows := sqlmock.
			NewRows([]string{"id", "product_id", "quantity", "updated_at"}).
			AddRow(1, 1, 10, dayCreated)
		mock.
			ExpectQuery(regexp.QuoteMeta("SELECT * FROM `product_quantity` WHERE product_id in (?) FOR UPDATE")).
			WithArgs(1).
			WillReturnRows(rows)

		// validate and update product quantity
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `product_quantity` SET `product_id`=?,`quantity`=?,`updated_at`=? WHERE `id` = ?")).
			WithArgs(1, 9, AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// checkout 1 of 10 existing items
		err := repo.SubmitCheckout(&entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:  &entity.Product{ID: 1, Serial: "120P90", Name: "Google Home", Price: 49.99, UpdatedAt: dayCreated},
					Quantity: 1,
				},
			},
		})
		assert.Nil(t, err)
	})

	t.Run("negative, item quantity is insufficient", func(t *testing.T) {
		mock.ExpectBegin()

		// lock for update product_quantity
		rows := sqlmock.
			NewRows([]string{"id", "product_id", "quantity", "updated_at"}).
			AddRow(1, 1, 10, dayCreated)
		mock.
			ExpectQuery(regexp.QuoteMeta("SELECT * FROM `product_quantity` WHERE product_id in (?) FOR UPDATE")).
			WithArgs(1).
			WillReturnRows(rows)

		// item is insufficient, rollback before update
		mock.ExpectRollback()

		// checkout 11 of 10 existing items
		err := repo.SubmitCheckout(&entity.Checkout{
			Items: []*entity.CheckoutItem{
				{
					Product:  &entity.Product{ID: 1, Serial: "120P90", Name: "Google Home", Price: 49.99, UpdatedAt: dayCreated},
					Quantity: 11,
				},
			},
		})
		assert.NotNil(t, err)
	})
}
