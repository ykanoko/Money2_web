package domain

type Money struct {
	ID    int32
	Name  string
	Price int64
	// Description string
	CategoryID int64
	UserID     int64
	// Image       []byte
	// Status      ItemStatus
	CreatedAt string
	UpdatedAt string
}

type Type struct {
	ID   int64
	Name string
}
