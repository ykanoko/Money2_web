package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/ykanoko/Money2_web/backend/db"
	"github.com/ykanoko/Money2_web/backend/handler"
)

// TODO:大量のユーザーに耐えられる仕組み作り
// DO:Dlete, Update機能追加

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(run())
}

func run() int {
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
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
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
	sqlDB, err := sql.Open("postgres", os.Getenv("CONN_STR"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer sqlDB.Close()

    var version string
    if err := sqlDB.QueryRow("select version()").Scan(&version); err != nil {
        panic(err)
    }
    fmt.Printf("version=%s\n", version)

	h := handler.Handler{
		DB:        sqlDB,
		UserRepo:  db.NewUserRepository(sqlDB),
		MoneyRepo: db.NewMoneyRepository(sqlDB),
	}

	// Routes
	// e.POST("/initialize", h.Initialize)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/log", h.AccessLog)

	// Login required
	l := e.Group("")
	l.Use(echojwt.WithConfig(config))

	// DO:groupでまとめる
	l.POST("/record_income", h.AddIncomeRecord)
	l.POST("/record_expense", h.AddExpenseRecord)
	l.POST("/record_pair_expense", h.AddPairExpenseRecord)
	l.DELETE("/record_income/:record_id", h.CancelIncomeRecord)

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
