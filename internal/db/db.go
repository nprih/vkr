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

	err = db.Ping()
	if err != nil {
		log.Println("ping db:", err)
		return nil, err
	}
	return db, nil
}

func Register(login string, password string) (message string, err error) {
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

func GetUserByLogin(login string) (*User, error) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	var user User
	err = conn.QueryRow("SELECT login, password, is_admin FROM users WHERE login = $1", login).
		Scan(&user.Login, &user.Password, &user.Is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с Login %s не найден", login)
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}
