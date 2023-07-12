package db

import (
	"context"
	"database/sql"

	"github.com/ykanoko/Money2_web/backend/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (int64, error)
	AddPair(ctx context.Context, pair domain.Pair) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	GetPair(ctx context.Context, id int64) (domain.Pair, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	UpdateBalance(ctx context.Context, id int64, balance float64) error
	UpdateCalculationUser1(ctx context.Context, id int64, calculation_user1 float64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

func (r *UserDBRepository) AddUser(ctx context.Context, user domain.User) (int64, error) {
	if _, err := r.ExecContext(ctx, "INSERT INTO users (name, balance) VALUES (?, 0)", user.Name); err != nil {
		return 0, err
	}
	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	row := r.QueryRowContext(ctx, "SELECT id FROM users WHERE rowid = LAST_INSERT_ROWID()")

	var id int64
	return id, row.Scan(&id)
}

func (r *UserDBRepository) AddPair(ctx context.Context, pair domain.Pair) (int64, error) {
	if _, err := r.ExecContext(ctx, "INSERT INTO pairs (password, user1_id, user2_id, calculation_user1) VALUES (?, ?, ?, 0)", pair.Password, pair.User1ID, pair.User2ID); err != nil {
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
	return user, row.Scan(&user.ID, &user.Name, &user.Balance)
}

func (r *UserDBRepository) GetPair(ctx context.Context, id int64) (domain.Pair, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM pairs WHERE id = ?", id)

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

func (r *UserDBRepository) UpdateBalance(ctx context.Context, id int64, balance float64) error {
	if _, err := r.ExecContext(ctx, "UPDATE users SET balance = ? WHERE id = ?", balance, id); err != nil {
		return err
	}
	return nil
}
func (r *UserDBRepository) UpdateCalculationUser1(ctx context.Context, id int64, calculation_user1 float64) error {
	if _, err := r.ExecContext(ctx, "UPDATE pairs SET calculation_user1 = ? WHERE id = ?", calculation_user1, id); err != nil {
		return err
	}
	return nil
}

type MoneyRepository interface {
	AddMoneyRecord(ctx context.Context, money domain.Money) (domain.Money, error)
	GetLatestMoneyRecordByPairID(ctx context.Context, pair_id int64) (domain.Money, error)
	GetMoneyRecordsByPairID(ctx context.Context, pair_id int64) ([]domain.Money, error)
	GetType(ctx context.Context, id int32) (domain.Type, error)
	GetTypes(ctx context.Context) ([]domain.Type, error)
}
type MoneyDBRepository struct {
	*sql.DB
}

func NewMoneyRepository(db *sql.DB) MoneyRepository {
	return &MoneyDBRepository{DB: db}
}

func (r *MoneyDBRepository) AddMoneyRecord(ctx context.Context, money domain.Money) (domain.Money, error) {
	if _, err := r.ExecContext(ctx, "INSERT INTO money2 (pair_id, type_id, user_id, amount) VALUES (?, ?, ?, ?)", money.PairID, money.TypeID, money.UserID, money.Amount); err != nil {
		return domain.Money{}, err
	}
	// TODO: if other insert query is executed at the same time, it might return wrong id
	// http.StatusConflict(409) 既に同じIDがあった場合
	row := r.QueryRowContext(ctx, "SELECT * FROM money2 WHERE rowid = LAST_INSERT_ROWID()")

	var res domain.Money
	return res, row.Scan(&res.ID, &res.PairID, &res.TypeID, &res.UserID, &res.Amount, &res.CreatedAt)
}

func (r *MoneyDBRepository) GetLatestMoneyRecordByPairID(ctx context.Context, pair_id int64) (domain.Money, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM money2 WHERE pair_id = ? ORDER BY created_at DESC LIMIT 1", pair_id)
	var money domain.Money
	return money, row.Scan(&money.ID, &money.PairID, &money.TypeID, &money.UserID, &money.Amount, &money.CreatedAt)
}

func (r *MoneyDBRepository) GetMoneyRecordsByPairID(ctx context.Context, pair_id int64) ([]domain.Money, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM money2 WHERE pair_id = ? ORDER BY created_at DESC LIMIT 30", pair_id)
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

func (r *MoneyDBRepository) GetType(ctx context.Context, id int32) (domain.Type, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM types WHERE id = ?", id)

	var typ domain.Type
	return typ, row.Scan(&typ.ID, &typ.Name)
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
