package handler

import (
	"database/sql"
	"fmt"
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
	PairID int64 `json:"pair_id"`
	jwt.RegisteredClaims
}

type InitializeResponse struct {
	Message string `json:"message"`
}

type registerRequest struct {
	User1Name string `json:"user1_name" validate:"required"`
	User2Name string `json:"user2_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type registerResponse struct {
	PairID    int64  `json:"pair_id"`
	User1ID   int64  `json:"user1_id"`
	User1Name string `json:"user1_name"`
	User2ID   int64  `json:"user2_id"`
	User2Name string `json:"user2_name"`
}

type loginRequest struct {
	PairID   int64  `json:"pair_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	PairID  int64  `json:"id"`
	User1ID int64  `json:"user1_id"`
	User2ID int64  `json:"user2_id"`
	Token   string `json:"token"`
}

type getPairStatusReponse struct {
	BalanceUser1 float64 `json:"balance_user1"`
	BalanceUser2 float64 `json:"balance_user2"`
	PayUser      string  `json:"pay_user"`
	PayAmount    float64 `json:"pay_amount"`
}

// type moneyRecordData struct {
type getMoneyRecordsResponse struct {
	Money2ID int64  `json:"money2_id"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	User     string `json:"user"`
	Amount   int64  `json:"amount"`
}

// type getMoneyRecordsResponse struct {
// 	Records []moneyRecordData `json:"records"`
// }

type addIncomeRecordRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
	Amount int64 `form:"amount" validate:"required"`
}

type addIncomeRecordRecordResponse struct {
	Money2ID int64 `json:"money2_id"`
}

type addPairExpenseRecordRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
	Amount int64 `form:"amount" validate:"required"`
}

type addPairExpenseRecordResponse struct {
	Money2ID int64 `json:"money2_id"`
}

type addIndivisualExpenseRecordRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
	Amount int64 `form:"amount" validate:"required"`
}

type addIndivisualExpenseRecordResponse struct {
	Money2ID int64 `json:"money2_id"`
}

type Handler struct {
	DB        *sql.DB
	UserRepo  db.UserRepository
	MoneyRepo db.MoneyRepository
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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
		return echo.NewHTTPError(http.StatusBadRequest, "name1, name2 and password are all required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user1ID, err := h.UserRepo.AddUser(c.Request().Context(), domain.User{Name: req.User1Name})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user2ID, err := h.UserRepo.AddUser(c.Request().Context(), domain.User{Name: req.User2Name})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	pairID, err := h.UserRepo.AddPair(c.Request().Context(), domain.Pair{Password: string(hash), User1ID: user1ID, User2ID: user2ID})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, registerResponse{PairID: pairID, User1ID: user1ID, User1Name: req.User1Name, User2ID: user2ID, User2Name: req.User2Name})
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

	pair, err := h.UserRepo.GetPair(ctx, req.PairID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pair.Password), []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		req.PairID,
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
		PairID:  pair.ID,
		User1ID: pair.User1ID,
		User2ID: pair.User2ID,
		Token:   encodedToken,
	})
}

func getPairID(c echo.Context) (int64, error) {
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return -1, fmt.Errorf("invalid token")
	}
	claims := user.Claims.(*JwtCustomClaims)
	if claims == nil {
		return -1, fmt.Errorf("invalid token")
	}
	return claims.PairID, nil
}

// // DO:Creat!!!!!!!
// func (h *Handler) AddIncomeRecord(c echo.Context) error {
// 	ctx := c.Request().Context()

// 	req := new(addMoneyRecordRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}
// }

// // DO:てーぶるの変更を直す（精算column）
// func (h *Handler) AddPairExpenseRecord(c echo.Context) error {
// 	ctx := c.Request().Context()

// 	req := new(addMoneyRecordRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "all columns are required")
// 	}

// 	_, err := h.MoneyRepo.GetType(ctx, req.TypeID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
// 		}
// 		return echo.NewHTTPError(http.StatusInternalServerError, err)
// 	}

// 	user, err := h.UserRepo.GetUser(ctx, req.UserID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}

// 	users, err := h.UserRepo.GetUsers(ctx)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}

// 	// latestMoneyRecord, err := h.MoneyRepo.GetLatestMoneyRecord(ctx)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusUnauthorized, err)
// 	}
// 	calculationAmount := float64(req.Amount) / 2

// 	// 残金の変更
// 	if req.TypeID == 1 {
// 		if err := h.UserRepo.UpdateBalance(ctx, req.UserID, user.Balance+float64(req.Amount)); err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
// 		}

// 	} else if req.TypeID == 2 {
// 		for _, use := range users {
// 			if err := h.UserRepo.UpdateBalance(ctx, use.ID, use.Balance-float64(req.Amount)/2); err != nil {
// 				return echo.NewHTTPError(http.StatusInternalServerError, err)
// 			}
// 		}
// 		// 精算の金額を変更
// 		if req.UserID == 2 {
// 			calculationAmount = -calculationAmount
// 		}
// 	}

// 	//money2 tableに登録
// 	moneyRecord, err := h.MoneyRepo.AddMoneyRecord(c.Request().Context(), domain.Money{
// 		TypeID: req.TypeID,
// 		UserID: req.UserID,
// 		Amount: req.Amount,
// 		// CalculationUser1: latestMoneyRecord.CalculationUser1 + calculationAmount,
// 	})
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err)
// 	}

// 	return c.JSON(http.StatusOK, addMoneyRecordResponse{ID: int64(moneyRecord.ID)})
// }

func (h *Handler) GetPairStatus(c echo.Context) error {
	ctx := c.Request().Context()

	pairID, err := getPairID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	pair, err := h.UserRepo.GetPair(ctx, pairID)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Pair not found.")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user1, err := h.UserRepo.GetUser(ctx, pair.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user2, err := h.UserRepo.GetUser(ctx, pair.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var payUser string
	if pair.CalculationUser1 < 0 {
		payUser = user1.Name
	} else if pair.CalculationUser1 > 0 {
		payUser = user2.Name
	} else {
		payUser = ""
	}

	return c.JSON(http.StatusOK, getPairStatusReponse{
		BalanceUser1: user1.Balance,
		BalanceUser2: user2.Balance,
		PayUser:      payUser,
		PayAmount:    math.Abs(pair.CalculationUser1),
	})
}

func (h *Handler) GetMoneyRecords(c echo.Context) error {
	ctx := c.Request().Context()
	pairID, err := getPairID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	pair, err := h.UserRepo.GetPair(ctx, pairID)
	// TODO: not found handling
	// http.StatusNotFound(404)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Pair not found.")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user1, err := h.UserRepo.GetUser(ctx, pair.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user2, err := h.UserRepo.GetUser(ctx, pair.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	moneyRecords, err := h.MoneyRepo.GetMoneyRecordsByPairID(ctx, pairID)
	// TODO: not found handling
	// http.StatusNotFound(404)
	// DO:ペアにおけるレコードが無い場合、Rsponseはnullになっている
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Record not found.")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	types, err := h.MoneyRepo.GetTypes(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var typeName string
	var userName string
	var res []getMoneyRecordsResponse
	for _, moneyRecord := range moneyRecords {
		for _, typ := range types {
			if typ.ID == moneyRecord.TypeID {
				typeName = typ.Name
			}
		}

		if user1.ID == moneyRecord.UserID {
			userName = user1.Name
		} else if user2.ID == moneyRecord.UserID {
			userName = user2.Name
		}
		res = append(res, getMoneyRecordsResponse{Money2ID: moneyRecord.ID, Date: moneyRecord.CreatedAt, Type: typeName, User: userName, Amount: moneyRecord.Amount})
	}

	// res := getMoneyRecordsResponse{Records: resMoneyRecords}
	return c.JSON(http.StatusOK, res)
}

// // DO:Creat!!!!!!!
// func (h *Handler) AddIndivisualExpenseRecord(c echo.Context) error {
// 	ctx := c.Request().Context()

// 	req := new(addMoneyRecordRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}
// }
