package db

import (
	"database/sql"
	"log"
	"vkr/internal/service"

	_ "modernc.org/sqlite"
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "volumes/photo_bank")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("ping db:", err)
		return nil, err
	}
	return db, nil
}

func RegisterUser(login string, password string) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	pass, err := service.HashPassword(password)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(login, pass)
}
