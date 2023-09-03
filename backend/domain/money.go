package domain

import "time"

type Money struct {
	ID int64
	// DO:ある程度記録が溜まったら、最初の方を削除？
	PairID int64
	TypeID int32
	// 収支を行うUserID
	UserID    int64
	Amount    int64
	CreatedAt time.Time
	// UpdatedAt        string DO:フロントエンドから修正できるようにする？
}

type Type struct {
	ID   int32
	Name string
}
