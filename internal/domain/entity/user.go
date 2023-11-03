package entity

type User struct {
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	ID        int64
	IP        string
	UserAgent string
}
