package main

import (
	"flag"
	"fmt"
	"hometest1/config"
	"hometest1/core/entity"
	"hometest1/core/module"
	"hometest1/handler"
	productrepository "hometest1/repository/product-repository"
	promotionrepository "hometest1/repository/promotion-repository"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var loadDotEnv = flag.Bool("loadDotEnv", false, "load .env file into ENV")

func main() {
	flag.Parse()

	// load .env file?
	if loadDotEnv != nil && *loadDotEnv {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err.Error())
		}
	}

	// load config
	cfg := config.Get()
	db := config.Connect()

	// load repository
	productRepo := productrepository.New(db)
	promoRepo := promotionrepository.New(db)

	// load usecase
	checkoutUC := module.NewCheckoutUsecase(productRepo, promoRepo)

	// load handler
	checkoutHandler := handler.NewCheckoutHandler(checkoutUC)

	// load echo framework
	e := echo.New()
	// set echo validator
	e.Validator = &CustomValidator{validator: validator.New()}
	// set error handler
	e.HTTPErrorHandler = errorHandler

	// route
	e.POST("/checkout", checkoutHandler.Submit)

	// run
	e.Logger.Fatal(e.Start(":" + cfg.HttpPort))
}

func errorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				report.Message = fmt.Sprintf("%s is required",
					err.Field())
			}
		}
	}

	if entityError, ok := err.(entity.Err); ok {
		report.Code = entityError.GetCode()
	}

	c.JSON(report.Code, report)
}
