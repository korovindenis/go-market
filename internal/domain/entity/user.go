package entity

type User struct {
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	ID        uint64
	IP        string
	UserAgent string
}
