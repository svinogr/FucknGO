package model

type UserModel struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
