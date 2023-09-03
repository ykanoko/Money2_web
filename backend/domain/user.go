package domain

import "time"

type Pair struct {
	ID       int64
	Password string
	User1ID  int64
	User2ID  int64
	// User1を基準とした精算額
	CalculationUser1 float64
	CreatedAt        time.Time
}

type User struct {
	ID      int64
	Name    string
	Balance float64
}
