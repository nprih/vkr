package db

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Is_admin bool   `json:"is_admin"`
}
