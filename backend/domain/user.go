package domain

type User struct {
	ID       int64
	Password string
	Name     string
	Balance  float64
	// Userを基準とした精算額
	Calculation float64
}
