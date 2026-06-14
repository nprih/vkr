package db

import "time"

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Is_admin bool   `json:"is_admin"`
}

type Image struct {
	Id           int64
	UserId       int64
	Filename     string
	OriginalName string
	FilePath     string
	FileSize     int
	MimeType     string
	Description  string
	CreatedAt    time.Time
}

type UserImage struct {
	Id       int64
	FilePath string
	Author   string
}
