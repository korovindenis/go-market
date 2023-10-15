package entity

type User struct {
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Ip        string `json:"-"`
	UserAgent string `json:"-"`
}
