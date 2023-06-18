package domain

type Money struct {
	ID int32
	// DO:ある程度記録が溜まったら、最初の方を削除？
	TypeID int64
	// 収支を行うUserID
	UserID int64
	Amount int64
	// User1を基準とした精算額
	CalculationUser1 int64
	CreatedAt        string
	// UpdatedAt        string DO:フロントエンドから修正できるようにする？
}

type Type struct {
	ID   int32
	Name string
}
