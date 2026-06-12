package db

import (
	"database/sql"
	"fmt"
	"log"

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

func RegisterUser(login string, password string) (message string, err error) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer conn.Close()

	_, err = conn.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("Пользователь %s добавлен", login), nil
}
