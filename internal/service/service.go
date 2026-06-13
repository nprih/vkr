package service

import (
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hash, nil
}

func CheckPassword(password string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем сессию
		session := sessions.Default(c)

		// Проверяем, есть ли пользователь в сессии
		login := session.Get("login")
		if login == nil {
			// Пользователь не авторизован - перенаправляем на логин
			c.Redirect(http.StatusFound, "/login")
			c.Abort() // Прерываем выполнение запроса
			return
		}

		// Пользователь авторизован - продолжаем
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем сессию
		session := sessions.Default(c)

		// Проверяем права администратора
		isAdmin, ok := session.Get("is_admin").(bool)
		if !ok || !isAdmin {
			// Нет прав администратора
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":   "Доступ запрещен",
				"message": "Требуются права администратора",
			})
			c.Abort()
			return
		}

		// Пользователь - администратор, продолжаем
		c.Next()
	}
}
