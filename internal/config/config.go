package config

import (
	"log"
	"os"
	"strconv"
	"vkr/internal/repository"

	"github.com/joho/godotenv"
)

var Port string
var SessionSecret string
var SessionLifetime string
var SessionName string
var DbPath string

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or error loading it")
		Port = "8085"
		SessionSecret = "vcFCnDZjK7ohNCba6cAwU89H7284MK8E"
		SessionLifetime = strconv.Itoa(60 * 60 * 1)
		SessionName = "userSession"
		DbPath = "./volumes/photo_bank"
		return
	}

	Port = os.Getenv("APP_PORT")
	SessionSecret = os.Getenv("SESSION_SECRET")
	SessionLifetime = os.Getenv("SESSION_LIFETIME")
	SessionName = os.Getenv("SESSION_NAME")
	DbPath = "./volumes/photo_bank"

	if err := repository.InitDB(DbPath); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
}
