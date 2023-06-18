package handler

import (
	"database/sql"
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

// type getMoney2Response struct {
// 	ID         int32  `json:"id"`
// 	Date       string `json:"date"`
// 	Type       string `json:"type"`
// 	User       string `json:"user"`
// 	Amount     int64  `json:"amount"`
// 	MoneyUser1 int64  `json:"money_user1"`
// 	MoneyUser2 int64  `json:"money_user2"`
// 	PayUser    string `json:"pay_user"`
// 	PayAmount  int64  `json:"pay_amount"`
// }

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

// func (h *Handler) GetMoney(c echo.Context) error {
// ctx := c.Request().Context()

// money2, err := h.MoneyRepo.GetMoney(ctx)
// // TODO: not found handling
// // http.StatusNotFound(404)
// if err != nil {
// 	if err == sql.ErrNoRows {
// 		return echo.NewHTTPError(http.StatusNotFound, "Record not found.")
// 	}
// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
// }

// var res []getMoney2Response
// for _, money := range money2 {
// 	if money.
// 	types, err := h.MoneyRepo.GetTypes(ctx)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}
// 	users, err := h.UserRepo.GetUsers(ctx)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}

// 	for _, typ := range types {
// 		if typ.ID == money.TypeID {
// 			for _, user:= range users {
// 				if user.ID == money.UserID {
// 					userName := user.Name
// 				}
// 				if user.ID ==  {
// 					userName := user.Name
// 				}
// 			}
// 			res = append(res, getMoney2Response{ID: money.ID, Date: money.Date, Type: Typ.Name, User: User.Name, Amount: money.Amount, MoneyUser1: money.MoneyUser1, MoneyUser2: money.MoneyUser2, PayUser: user.Name, PayAmount: math.Abs(float64(money.CalculationUser1))})

// 		}
// 	}
// }

// return c.JSON(http.StatusOK, res)
// }

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
