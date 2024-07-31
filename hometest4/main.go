package main

import (
	"flag"
	"fmt"
	"hometest4/core/entity"
	"hometest4/core/handler"
	customMiddleware "hometest4/core/middleware"
	"hometest4/core/repository"
	"hometest4/core/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var port = flag.Int("port", 8080, "http port")

func main() {
	flag.Parse()

	// init data
	user := make(map[string]*entity.User)
	topup := make(map[string][]*entity.Transaction)
	transfer := make(map[string][]*entity.Transaction)

	// init repo
	repo := repository.NewRepository(user, topup, transfer)

	// init service
	srv := service.NewService(repo)

	// init handler
	h := handler.NewHandler(srv)

	// init echo
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	middlewares := []echo.MiddlewareFunc{
		customMiddleware.GetJWT(),
		customMiddleware.GetUserFromJWT(),
	}

	// error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		herr, ok := err.(*echo.HTTPError)
		if !ok {
			herr = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if herr.Code == http.StatusUnauthorized {
			herr.Message = "User token not valid"
		}
		c.String(herr.Code, fmt.Sprintf("%v", herr.Message))
	}

	// router
	e.GET("/balance_read", h.BalanceRead, middlewares...)
	e.POST("/transfer", h.Transfer, middlewares...)
	e.GET("/top_users", h.ListOverallTopTrxUsersByValue, middlewares...)
	e.POST("/create_user", h.RegisterUser)
	e.POST("/balance_topup", h.BalanceTopup, middlewares...)
	e.GET("/top_transactions_per_user", h.TopTrxForUser, middlewares...)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
