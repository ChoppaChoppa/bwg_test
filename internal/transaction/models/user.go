package models

type User struct {
	ID       int    `json:"ID" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type Balance struct {
	ID      int     `json:"ID" db:"id"`
	UserID  int     `json:"userID" db:"user_id"`
	Balance float64 `json:"balance" db:"balance"`
}
