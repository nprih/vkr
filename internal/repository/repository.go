package repository

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func InitDB(dataSourceName string) error {
	var err error
	once.Do(func() {
		dbInstance, err = sql.Open("sqlite3", dataSourceName)
		if err != nil {
			log.Printf("Ошибка открытия БД: %v", err)
			return
		}

		dbInstance.SetMaxOpenConns(1)
		dbInstance.SetMaxIdleConns(1)

		if err = dbInstance.Ping(); err != nil {
			log.Printf("Ошибка пинга БД: %v", err)
			return
		}
		log.Println("База данных успешно инициализирована.")
	})
	return err
}

func GetDB() *sql.DB {
	if dbInstance == nil {
		log.Fatal("База данных не инициализирована.")
	}
	return dbInstance
}

func CloseDB() error {
	if dbInstance != nil {
		return dbInstance.Close()
	}
	return nil
}
