package model

type User struct {
	Id       string `json:"id",omitempty`
	Username string `json:"username"`
	Password string `json:"password"`
}
