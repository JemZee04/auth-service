package model

type User struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}
