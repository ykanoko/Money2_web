package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ykanoko/Money2_web/backend/db"
	"github.com/ykanoko/Money2_web/backend/handler"
)

/*
複数行の
コメント
*/

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	logfile := os.Getenv("LOGFILE")
	if logfile == "" {
		logfile = "access.log"
	}
	lf, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: logFormat(),
		Output: io.MultiWriter(os.Stdout, lf),
	})
	e.Use(logger)
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("failed to load environment variables", err)
	}
	frontURL := os.Getenv("FRONT_URL")
	fmt.Println("frontURL", frontURL)
	if frontURL == "" {
		// frontURL = "http://localhost:3000"
		frontURL = "https://money2-web.vercel.app/"
	}
	fmt.Println("frontURLCORS", []string{frontURL})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "DELETE", "OPTIONS", "POST"},
	}))
	e.Use(middleware.BodyLimit("5M"))

	// jwt DO:分からないで使うと脆弱性あり。勉強する。
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte(handler.GetSecret()),
	}

	// db
	sqlDB, err := db.PrepareDB(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare DB: %s\n", err)
		return exitError
	}
	defer sqlDB.Close()

	h := handler.Handler{
		DB:        sqlDB,
		UserRepo:  db.NewUserRepository(sqlDB),
		MoneyRepo: db.NewMoneyRepository(sqlDB),
	}

	// Routes
	e.POST("/initialize", h.Initialize)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/log", h.AccessLog)

	// Login required
	l := e.Group("")
	l.Use(echojwt.WithConfig(config))

	l.POST("/record_income", h.AddIncomeRecord)
	l.POST("/record_pair_expense", h.AddPairExpenseRecord)
	// l.POST("/record_indivisual_expense", h.AddIndivisualExpenseRecord)
	l.GET("/pair_status", h.GetPairStatus)
	l.GET("/money_records", h.GetMoneyRecords)

	// Start server
	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return exitOK
}

func logFormat() string {
	// Customize freely: https://echo.labstack.com/guide/customization/
	var format string
	format += "time:${time_rfc3339}\t"
	format += "status:${status}\t"
	format += "method:${method}\t"
	format += "uri:${uri}\t"
	format += "latency:${latency_human}\t"
	format += "error:${error}\t"
	format += "\n"

	// Other log choice
	// - json format
	// `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
	// 	`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
	// 	`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
	// 	`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
	// - structured logging:  https://pkg.go.dev/golang.org/x/exp/slog

	return format
}
