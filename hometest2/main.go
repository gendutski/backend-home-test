package main

import (
	"context"
	"flag"
	"fmt"
	"hometest2/module"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	maxRange     int = 100
	maxGoroutine int = 1000
)

var port = flag.String("port", "8080", "set http port")
var simulateTimeout = flag.Int("simulateTimeout", 0, "set sleep in miliseconds, to check if service can handle timeout")

type Parser struct {
	Value  int
	Result string
}

func main() {
	// parse flags
	flag.Parse()

	// init echo
	e := echo.New()
	// custom middleware log
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// init
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Create a response recorder
			responseWriter := res.Writer
			recorder := &responseRecorder{responseWriter, http.StatusOK, []byte{}}
			res.Writer = recorder

			// Call next handler
			if err := next(c); err != nil {
				c.Error(err)
			}

			// set log
			var result string
			result += fmt.Sprintf("\n\tStatus: %d", recorder.statusCode)
			result += fmt.Sprintf("\n\tMethod: %s", req.Method)
			// endpoint
			endPoint := req.URL.Path
			if endPoint == "" {
				endPoint = "/"
			}
			result += fmt.Sprintf("\n\tEndpoint: %s", endPoint)
			// request
			result += fmt.Sprintf("\n\tRequest: %v", req.URL.Query())
			// response
			result += fmt.Sprintf("\n\tResponse: %s", string(recorder.body))
			// latency
			result += fmt.Sprintf("\n\tLatency: %s", time.Since(start))

			// print log
			log.Println(result)

			return nil
		}
	})
	// timeout midleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: "request timeout",
		Timeout:      time.Second,
	}))

	// route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/range-fizzbuzz", handleRangeFizzBuzz)

	// Start server
	go func() {
		if err := e.Start(":" + *port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("cannot run service at port: `%v`", *port)
		}
	}()

	// Channel to listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Block until a signal is received
	<-stop

	fmt.Println()
	log.Println("Shutting down the server...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Could not gracefully shutdown the server:", err)
	}

	log.Println("Server stopped")
}

func handleRangeFizzBuzz(c echo.Context) error {
	// get query params
	from, to, err := validateRangeFizzBuzz(c)
	if err != nil {
		return err
	}

	// set buffer
	ch := make(chan Parser, maxGoroutine)

	// send data
	for i := from; i <= to; i++ {
		go func(i int) {
			ch <- Parser{
				Value:  i,
				Result: module.SingleFizzBuzz(i),
			}
		}(i)
	}

	var parser []Parser
	// get data
	for i := from; i <= to; i++ {
		select {
		case <-c.Request().Context().Done():
			close(ch)
			// break loop
			i = to
		case p, ok := <-ch:
			if ok {
				parser = append(parser, p)
			}
		}

		if *simulateTimeout > 0 {
			time.Sleep(time.Millisecond * time.Duration(*simulateTimeout))
		}
	}
	// sort data
	sort.Slice(parser, func(i, j int) bool {
		return parser[i].Value < parser[j].Value
	})
	// set result
	var result string
	for i := range parser {
		if i != 0 {
			result += " "
		}
		result += parser[i].Result
	}

	return c.String(http.StatusOK, result)
}

func validateRangeFizzBuzz(c echo.Context) (from, to int, err error) {
	// validate from
	if _from := c.QueryParam("from"); _from != "" {
		from, err = strconv.Atoi(_from)
		if err != nil {
			err = echo.NewHTTPError(http.StatusBadRequest, "invalid `form` parameter, must be number")
			return
		}
	}

	// validate to
	if _to := c.QueryParam("to"); _to != "" {
		to, err = strconv.Atoi(_to)
		if err != nil {
			err = echo.NewHTTPError(http.StatusBadRequest, "invalid `to` parameter, must be number")
			return
		}
	}

	// validate range
	if from > to {
		err = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("`from` parameter must less or equal than %d", to))
		return
	}
	if to-from+1 > maxRange {
		err = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("range between `from` and `to` parameters must less than or equal tp %d", maxRange))
	}
	return
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
