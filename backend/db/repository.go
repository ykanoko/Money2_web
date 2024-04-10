package db

import (
	"context"
	"database/sql"

	"github.com/ykanoko/Money2_web/backend/domain"
)

type UserRepository interface {
	AddUser(tx *sql.Tx, user domain.User) (int64, error)
	AddPair(tx *sql.Tx, pair domain.Pair) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	GetPair(ctx context.Context, id int64) (domain.Pair, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	UpdateBalance(tx *sql.Tx, id int64, balance float64) error
	UpdateCalculationUser1(tx *sql.Tx, id int64, calculation_user1 float64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

func (r *UserDBRepository) AddUser(tx *sql.Tx, user domain.User) (int64, error) {
	row := tx.QueryRow("INSERT INTO users (name, balance) VALUES ($1, 0) RETURNING id", user.Name)
	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) AddPair(tx *sql.Tx, pair domain.Pair) (int64, error) {
	row := tx.QueryRow("INSERT INTO pairs (password, user1_id, user2_id, calculation_user1) VALUES ($1, $2, $3, 0) RETURNING id", pair.Password, pair.User1ID, pair.User2ID)

	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Balance)
}

func (r *UserDBRepository) GetPair(ctx context.Context, id int64) (domain.Pair, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM pairs WHERE id = $1", id)

	var pair domain.Pair
	return pair, row.Scan(&pair.ID, &pair.Password, &pair.User1ID, &pair.User2ID, &pair.CalculationUser1, &pair.CreatedAt)
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
		if err := rows.Scan(&user.ID, &user.Name, &user.Balance); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserDBRepository) UpdateBalance(tx *sql.Tx, id int64, balance float64) error {
	if _, err := tx.Exec("UPDATE users SET balance = $1 WHERE id = $2", balance, id); err != nil {
		return err
	}
	return nil
}
func (r *UserDBRepository) UpdateCalculationUser1(tx *sql.Tx, id int64, calculation_user1 float64) error {
	if _, err := tx.Exec("UPDATE pairs SET calculation_user1 = $1 WHERE id = $2", calculation_user1, id); err != nil {
		return err
	}
	return nil
}

type MoneyRepository interface {
	AddMoneyRecord(tx *sql.Tx, money domain.Money) (domain.Money, error)
	DeleteMoneyRecordByID(tx *sql.Tx, id int64) error
	GetMoneyRecordByID(ctx context.Context, id int64) (domain.Money, error)
	GetLatestMoneyRecordByPairID(ctx context.Context, pair_id int64) (domain.Money, error)
	GetMoneyRecordsByPairID(ctx context.Context, pair_id int64) ([]domain.Money, error)
	GetTypeNameByID(ctx context.Context, id int32) (string, error)
	GetTypes(ctx context.Context) ([]domain.Type, error)
}
type MoneyDBRepository struct {
	*sql.DB
}

func NewMoneyRepository(db *sql.DB) MoneyRepository {
	return &MoneyDBRepository{DB: db}
}

func (r *MoneyDBRepository) AddMoneyRecord(tx *sql.Tx, money domain.Money) (domain.Money, error) {
	row := tx.QueryRow("INSERT INTO money2 (pair_id, type_id, user_id, amount) VALUES ($1, $2, $3, $4) RETURNING id, pair_id, type_id, user_id, amount, created_at", money.PairID, money.TypeID, money.UserID, money.Amount)
	var res domain.Money
	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	return res, row.Scan(&res.ID, &res.PairID, &res.TypeID, &res.UserID, &res.Amount, &res.CreatedAt)
}

func (r *MoneyDBRepository) DeleteMoneyRecordByID(tx *sql.Tx, id int64) error {
	if _, err := tx.Exec("DELETE FROM money2 WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func (r *MoneyDBRepository) GetMoneyRecordByID(ctx context.Context, id int64) (domain.Money, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM money2 WHERE id = $1", id)
	var money domain.Money
	return money, row.Scan(&money.ID, &money.PairID, &money.TypeID, &money.UserID, &money.Amount, &money.CreatedAt)
}

func (r *MoneyDBRepository) GetLatestMoneyRecordByPairID(ctx context.Context, pair_id int64) (domain.Money, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM money2 WHERE pair_id = $1 ORDER BY created_at DESC LIMIT 1", pair_id)
	var money domain.Money
	return money, row.Scan(&money.ID, &money.PairID, &money.TypeID, &money.UserID, &money.Amount, &money.CreatedAt)
}

func (r *MoneyDBRepository) GetMoneyRecordsByPairID(ctx context.Context, pair_id int64) ([]domain.Money, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM money2 WHERE pair_id = $1 ORDER BY created_at DESC LIMIT 100", pair_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var money2 []domain.Money
	for rows.Next() {
		var money domain.Money
		if err := rows.Scan(&money.ID, &money.PairID, &money.TypeID, &money.UserID, &money.Amount, &money.CreatedAt); err != nil {
			return nil, err
		}
		money2 = append(money2, money)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return money2, nil
}

func (r *MoneyDBRepository) GetTypeNameByID(ctx context.Context, id int32) (string, error) {
	row := r.QueryRowContext(ctx, "SELECT name FROM types WHERE id = $1", id)
	var name string
	return name, row.Scan(&name)
}

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
