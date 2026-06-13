package service

import (
	"log"
	"net/http"
	"vkr/internal/db"

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

func SetSession(c *gin.Context, user *db.User) {
	session := sessions.Default(c)
	session.Set("login", user.Login)
	session.Set("is_login", true)
	session.Set("is_admin", user.Is_admin)
	session.Save()
}

func UnsetSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("login")
	session.Delete("is_login")
	session.Delete("is_admin")
	session.Save()
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		login := session.Get("login")
		if login == nil {
			log.Println("Пользователь не авторизован, редирект на /login")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAdmin, ok := session.Get("is_admin").(bool)
		if !ok || !isAdmin {
			log.Println(" Админ не авторизован, редирект на /")
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":    "Доступ запрещен",
				"message":  "Требуются права администратора",
				"redirect": "/",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
