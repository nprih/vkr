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
		session := sessions.Default(c)
		user := session.Get("is_login")
		if user == nil {
			c.Redirect(http.StatusFound, "/")
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
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error": "Доступ запрещен",
			})
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
