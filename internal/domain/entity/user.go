package entity

type User struct {
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Id        uint64
	Ip        string
	UserAgent string
}
