package db

import (
	"context"
	"database/sql"

	"github.com/ykanoko/Money2_web/backend/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	// UpdateBalance(ctx context.Context, id int64, balance int64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

func (r *UserDBRepository) AddUser(ctx context.Context, user domain.User) (int64, error) {
	if _, err := r.ExecContext(ctx, "INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password); err != nil {
		return 0, err
	}
	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	row := r.QueryRowContext(ctx, "SELECT id FROM users WHERE rowid = LAST_INSERT_ROWID()")

	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Password)
}
func (r *UserDBRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

type MoneyRepository interface {
	// AddMoney(ctx context.Context, money domain.Money) (domain.Money, error)
	// GetMoney(ctx context.Context, id int32) (domain.Money, error)
	// GetMoneyImage(ctx context.Context, id int32) ([]byte, error)
	// GetMoney(ctx context.Context) ([]domain.Money, error)
	// GetMoney2ByUserID(ctx context.Context, userID int64) ([]domain.Money, error)
	// GetType(ctx context.Context, id int64) (domain.Type, error)
	GetTypes(ctx context.Context) ([]domain.Type, error)
	// UpdateMoneyStatus(ctx context.Context, id int32, status domain.MoneyStatus) error
}
type MoneyDBRepository struct {
	*sql.DB
}

func NewMoneyRepository(db *sql.DB) MoneyRepository {
	return &MoneyDBRepository{DB: db}
}

// func (r *MoneyDBRepository) AddMoney(ctx context.Context, money domain.Money) (domain.Money, error) {
// 	if _, err := r.ExecContext(ctx, "INSERT INTO money2 (name, price, category_id, seller_id) VALUES (?, ?, ?, ?)", money.Name, money.Price, money.CategoryID, money.UserID); err != nil {
// 		return domain.Money{}, err
// 	}
// 	// TODO: if other insert query is executed at the same time, it might return wrong id
// 	// http.StatusConflict(409) 既に同じIDがあった場合
// 	row := r.QueryRowContext(ctx, "SELECT * FROM money2 WHERE rowid = LAST_INSERT_ROWID()")

// 	var res domain.Money
// 	return res, row.Scan(&res.ID, &res.Name, &res.Price, &res.CategoryID, &res.UserID, &res.CreatedAt, &res.UpdatedAt)
// }

// func (r *MoneyDBRepository) GetMoney(ctx context.Context) ([]domain.Money, error) {
// 	rows, err := r.QueryContext(ctx, "SELECT * FROM money2 ORDER BY time DESC LIMIT 3")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var money2 []domain.Money
// 	for rows.Next() {
// 		var money domain.Money
// 		if err := rows.Scan(&money.ID, &money.Date, &money.TypeID, &money.UserID, &money.Amount, &money.MoneyUser1, &money.MoneyUser2, &money.CalculationUser1); err != nil {
// 			return nil, err
// 		}
// 		money2 = append(money2, money)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return money2, nil
// }

// func (r *MoneyDBRepository) GetMoney2ByUserID(ctx context.Context, userID int64) ([]domain.Money, error) {
// 	rows, err := r.QueryContext(ctx, "SELECT * FROM money2 WHERE seller_id = ?", userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var money2 []domain.Money
// 	for rows.Next() {
// 		var item domain.Money
// 		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		money2 = append(money2, item)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return money2, nil
// }

// func (r *MoneyDBRepository) GetType(ctx context.Context, id int64) (domain.Type, error) {
// 	row := r.QueryRowContext(ctx, "SELECT * FROM types WHERE id = ?", id)

// 	var cat domain.Type
// 	return cat, row.Scan(&cat.ID, &cat.Name)
// }

func (r *MoneyDBRepository) GetTypes(ctx context.Context) ([]domain.Type, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []domain.Type
	for rows.Next() {
		var typ domain.Type
		if err := rows.Scan(&typ.ID, &typ.Name); err != nil {
			return nil, err
		}
		types = append(types, typ)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return types, nil
}
