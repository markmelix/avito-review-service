package domain

type User struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
