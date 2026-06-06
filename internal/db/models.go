package db

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Pass     string `json:"pass"`
	Is_admin bool   `json:"is_admin"`
}
