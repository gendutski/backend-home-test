package promotionrepository_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"hometest1/core/entity"
	"hometest1/core/repository"
	promotionrepository "hometest1/repository/promotion-repository"

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

func initRepo(db *sql.DB, mock sqlmock.Sqlmock) (repository.PromotionRepo, error) {
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
	return promotionrepository.New(gdb), nil
}

func Test_GetPromotionByProducts(t *testing.T) {
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
			NewRows([]string{"id", "type", "product_id", "match_quantity", "promo_value", "promo_product_id", "updated_at", "deleted_at"}).
			AddRow(1, 1, 2, 1, 1, 4, dayCreated, nil).
			AddRow(3, 3, 3, 3, 10, 0, dayCreated, nil)

		mock.
			ExpectQuery(regexp.QuoteMeta("SELECT * FROM `promotion` WHERE product_id in (?,?) AND `promotion`.`deleted_at` IS NULL ORDER BY product_id asc, type asc")).
			WithArgs(int64(2), int64(3)).
			WillReturnRows(rows)

		resp, err := repo.GetPromotionByProducts([]*entity.Product{
			{ID: 2, Serial: "43N23P", Name: "MacBook Pro", Price: 5399.99, UpdatedAt: dayCreated},
			{ID: 3, Serial: "A304SD", Name: "Alexa Speaker", Price: 49.99, UpdatedAt: dayCreated},
		})
		assert.Nil(t, err)
		assert.Equal(t, map[int64][]*entity.Promotion{
			2: {{ID: 1, Type: 1, ProductID: 2, MatchQuantity: 1, PromoValue: 1, PromoProductID: 4, UpdatedAt: dayCreated}},
			3: {{ID: 3, Type: 3, ProductID: 3, MatchQuantity: 3, PromoValue: 10, PromoProductID: 0, UpdatedAt: dayCreated}},
		}, resp)
	})
}
