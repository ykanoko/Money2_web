package handler

import (
	"database/sql"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/ykanoko/Money2_web/backend/db"
	"github.com/ykanoko/Money2_web/backend/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	logFile = getEnv("LOGFILE", "access.log")
)

type JwtCustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type InitializeResponse struct {
	Message string `json:"message"`
}

type registerRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type registerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type loginRequest struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type getMoneyRecordsResponse struct {
	ID           int32   `json:"id"`
	Date         string  `json:"date"`
	Type         string  `json:"type"`
	User         string  `json:"user"`
	Amount       int64   `json:"amount"`
	BalanceUser1 float64 `json:"balance_user1"`
	BalanceUser2 float64 `json:"balance_user2"`
	PayUser      string  `json:"pay_user"`
	PayAmount    float64 `json:"pay_amount"`
}

type addMoneyRecordRequest struct {
	TypeID int32 `form:"type_id" validate:"required"`
	UserID int64 `form:"user_id" validate:"required"`
	Amount int64 `form:"amount" validate:"required"`
}

type addMoneyRecordResponse struct {
	ID int64 `json:"id"`
}

type Handler struct {
	DB        *sql.DB
	UserRepo  db.UserRepository
	MoneyRepo db.MoneyRepository
}

func GetSecret() string {
	if secret := os.Getenv("SECRET"); secret != "" {
		return secret
	}
	return "secret-key"
}

func (h *Handler) Initialize(c echo.Context) error {
	err := os.Truncate(logFile, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to truncate access log"))
	}

	err = db.Initialize(c.Request().Context(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to initialize"))
	}

	return c.JSON(http.StatusOK, InitializeResponse{Message: "Success"})
}

func (h *Handler) AccessLog(c echo.Context) error {
	return c.File(logFile)
}

func (h *Handler) Register(c echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "name and password are both required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	userID, err := h.UserRepo.AddUser(c.Request().Context(), domain.User{Name: req.Name, Password: string(hash)})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, registerResponse{ID: userID, Name: req.Name})
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id and password are both required")
	}

	user, err := h.UserRepo.GetUser(ctx, req.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		req.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	encodedToken, err := token.SignedString([]byte(GetSecret()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, loginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: encodedToken,
	})
}

func (h *Handler) GetMoneyRecords(c echo.Context) error {
	ctx := c.Request().Context()

	moneyRecords, err := h.MoneyRepo.GetMoneyRecords(ctx)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Record not found.")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []getMoneyRecordsResponse
	types, err := h.MoneyRepo.GetTypes(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	users, err := h.UserRepo.GetUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var userName string
	var balance1 float64
	var balance2 float64
	var payUserID int64
	var payUser string

	for _, moneyRecord := range moneyRecords {
		if moneyRecord.CalculationUser1 > 0 {
			payUserID = 1
		} else if moneyRecord.CalculationUser1 < 0 {
			payUserID = 2
		} else {
			payUserID = 0
		}

		for _, typ := range types {
			if typ.ID == moneyRecord.TypeID {
				for _, user := range users {
					if user.ID == 1 {
						balance1 = user.Balance

					} else if user.ID == 2 {
						balance2 = user.Balance
					}

					if user.ID == moneyRecord.UserID {
						userName = user.Name
					}

					// CalculationUser1 == 0　の時
					if payUserID == 0 {
						payUser = ""
					}
					if user.ID == payUserID {
						payUser = user.Name
					}
				}
				res = append(res, getMoneyRecordsResponse{ID: moneyRecord.ID, Date: moneyRecord.CreatedAt, Type: typ.Name, User: userName, Amount: moneyRecord.Amount, BalanceUser1: balance1, BalanceUser2: balance2, PayUser: payUser, PayAmount: math.Abs(moneyRecord.CalculationUser1)})
			}
		}
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AddMoneyRecord(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addMoneyRecordRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if req.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "add plus amount")
	}
	// DO:フロントエンド実装

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "all columns are required")
	}

	_, err := h.MoneyRepo.GetType(ctx, req.TypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	users, err := h.UserRepo.GetUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, user := range users {
		if err := h.UserRepo.UpdateBalance(ctx, user.ID, user.Balance-float64(req.Amount)/2); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	latestMoneyRecord, err := h.MoneyRepo.GetLatestMoneyRecord(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	calculationAmount := float64(req.Amount) / 2
	if req.UserID == 2 {
		calculationAmount = -calculationAmount
	}

	moneyRecord, err := h.MoneyRepo.AddMoneyRecord(c.Request().Context(), domain.Money{
		TypeID:           req.TypeID,
		UserID:           req.UserID,
		Amount:           req.Amount,
		CalculationUser1: latestMoneyRecord.CalculationUser1 + calculationAmount,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, addMoneyRecordResponse{ID: int64(moneyRecord.ID)})
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
