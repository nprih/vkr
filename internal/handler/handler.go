package handler

import (
	"log"
	"net/http"
	"vkr/internal/db"

	"github.com/gin-gonic/gin"
)

var req Credentials

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{
		"title": "Main website",
	})
}

func GetLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login website",
	})
}

func PostLoginHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("New login try:", req.Login)
}

func GetRegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Main website",
	})
}

func PostRegisterHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("New registry try:", req.Login, req.Password)
	db.RegisterUser(req.Login, req.Password)
}
