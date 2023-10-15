package entity

type User struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDevice struct {
	Ip        string
	UserAgent string
}
