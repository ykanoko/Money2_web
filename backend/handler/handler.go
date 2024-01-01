package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ykanoko/Money2_web/backend/db"
	"github.com/ykanoko/Money2_web/backend/domain"
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

type userReponse struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type registerRequest struct {
	User1Name string `json:"user1_name" validate:"required"`
	User2Name string `json:"user2_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type registerResponse struct {
	PairID int64       `json:"pair_id"`
	User1  userReponse `json:"user1"`
	User2  userReponse `json:"user2"`
}

type loginRequest struct {
	PairID   int64  `json:"pair_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	PairID int64       `json:"id"`
	User1  userReponse `json:"user1"`
	User2  userReponse `json:"user2"`
	Token  string      `json:"token"`
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

// linebotに送るメッセージ
type lineMessage struct {
	To       string `json:"to"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// DO:addはまとめる？
type addIncomeRecordRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
	Amount int64 `json:"amount" validate:"required"`
}

type addIncomeRecordResponse struct {
	CreatedAt time.Time `json:"created_at"`
}

type addExpenseRecordRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
	Amount int64 `json:"amount" validate:"required"`
}

type addExpenseRecordResponse struct {
	CreatedAt time.Time `json:"created_at"`
}

type addPairExpenseRecordRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
	Amount int64 `json:"amount" validate:"required"`
}

type addPairExpenseRecordResponse struct {
	CreatedAt time.Time `json:"created_at"`
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
	fmt.Println(req)
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

	return c.JSON(http.StatusOK, registerResponse{
		PairID: pairID,
		User1:  userReponse{ID: user1ID, Name: req.User1Name, Balance: 0},
		User2:  userReponse{ID: user2ID, Name: req.User2Name, Balance: 0},
	})
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Pair ID and password are both required")
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

	user1, err := h.UserRepo.GetUser(ctx, pair.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user2, err := h.UserRepo.GetUser(ctx, pair.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, loginResponse{
		PairID: pair.ID,
		User1:  userReponse{ID: user1.ID, Name: user1.Name, Balance: user1.Balance},
		User2:  userReponse{ID: user2.ID, Name: user2.Name, Balance: user2.Balance},
		Token:  encodedToken,
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

func (h *Handler) sendLineMessage(lineGroupID string, moneyRecord domain.Money, userName string) error {
	// linebotに送るメッセージ
	message := lineMessage{
		To: lineGroupID,
		Messages: []Message{
			{
				Type: "text",
				// DO:種類の部分の汎用性を上げる
				// DO:日付の表記の部分を関数に切り出す？
				Text: fmt.Sprintf("日付：%s\n種類：収入\n名前：%s\n金額：%d円", moneyRecord.CreatedAt.In(time.FixedZone("JST", 9*60*60)).Format("2006/01/02 15:04:05"), userName, moneyRecord.Amount),
			},
		},
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// linebotに送信
	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(messageJson))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	return nil
}	

func (h *Handler) AddIncomeRecord(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addIncomeRecordRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "all columns are required")
	}

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

	if req.UserID != pair.User1ID && req.UserID != pair.User2ID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "You don't belong to this pair.")
	}

	// 残金の変更
	user, err := h.UserRepo.GetUser(ctx, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := h.UserRepo.UpdateBalance(ctx, req.UserID, user.Balance+float64(req.Amount)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//money2 tableに登録
	moneyRecord, err := h.MoneyRepo.AddMoneyRecord(c.Request().Context(), domain.Money{
		PairID: pairID,
		TypeID: 1,
		// DO:typeTableいらない？
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// DO:groupIdを変数にする、環境変数にする
	// if err := h.sendLineMessage(os.Getenv("LINE_GROUP_ID"), moneyRecord, user.Name); err != nil {
	if err := h.sendLineMessage("Ca52523506465781a55cd4665945d7976", moneyRecord, user.Name); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// DO:ID返す意味ない？
	return c.JSON(http.StatusOK, addIncomeRecordResponse{CreatedAt: moneyRecord.CreatedAt})
}

func (h *Handler) AddExpenseRecord(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addExpenseRecordRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "all columns are required")
	}

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

	if req.UserID != pair.User1ID && req.UserID != pair.User2ID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "You don't belong to this pair.")
	}

	// 残金の変更
	user, err := h.UserRepo.GetUser(ctx, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := h.UserRepo.UpdateBalance(ctx, req.UserID, user.Balance-float64(req.Amount)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//money2 tableに登録
	moneyRecord, err := h.MoneyRepo.AddMoneyRecord(c.Request().Context(), domain.Money{
		PairID: pairID,
		TypeID: 2,
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, addExpenseRecordResponse{CreatedAt: moneyRecord.CreatedAt})
}

func (h *Handler) AddPairExpenseRecord(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addPairExpenseRecordRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "all columns are required")
	}

	if req.Amount <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Amount have to be >0")
	}

	// DO:関数に切り出したい?
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

	if req.UserID != pair.User1ID && req.UserID != pair.User2ID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "You don't belong to this pair.")
	}

	user1, err := h.UserRepo.GetUser(ctx, pair.User1ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user2, err := h.UserRepo.GetUser(ctx, pair.User2ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 残金の変更
	if err := h.UserRepo.UpdateBalance(ctx, user1.ID, user1.Balance-float64(req.Amount)/2); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := h.UserRepo.UpdateBalance(ctx, user2.ID, user2.Balance-float64(req.Amount)/2); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// 精算の金額を変更
	calculationAmount := float64(req.Amount) / 2
	if req.UserID == user2.ID {
		calculationAmount = -calculationAmount
	}
	if err := h.UserRepo.UpdateCalculationUser1(ctx, pairID, pair.CalculationUser1+calculationAmount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//money2 tableに登録
	moneyRecord, err := h.MoneyRepo.AddMoneyRecord(c.Request().Context(), domain.Money{
		PairID: pairID,
		TypeID: 3,
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, addPairExpenseRecordResponse{CreatedAt: moneyRecord.CreatedAt})
}

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

	// DO:388行まで、外に関数切り出す？
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

	// DO:引数1つ追加（取ってくる行数がフロントエンドで複数パターン必要）
	moneyRecords, err := h.MoneyRepo.GetMoneyRecordsByPairID(ctx, pairID)
	// TODO: not found handling
	// http.StatusNotFound(404)
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
	res := []getMoneyRecordsResponse{}
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
		res = append(res, getMoneyRecordsResponse{Money2ID: moneyRecord.ID, Date: moneyRecord.CreatedAt.In(time.FixedZone("JST", 9*60*60)).Format("2006/01/02 15:04:05"), Type: typeName, User: userName, Amount: moneyRecord.Amount})
	}
	return c.JSON(http.StatusOK, res)
}
