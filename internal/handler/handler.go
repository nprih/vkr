package handler

import (
	"log"
	"net/http"
	"vkr/internal/db"
	"vkr/internal/service"

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
	log.Println("New login try:", req.Login, req.Password)
	user, err := db.GetUserByLogin(req.Login)
	if err != nil {
		log.Printf("Пользователь %s не найден", req.Login)
		return
	}
	if service.CheckPassword(req.Password, user.Password) {
		log.Println("Пароль совпадает")
	} else {
		log.Println("Не верный пароль")
	}
}

func GetRegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Main website",
	})
}

func PostRegisterHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Не удалось прочитать отправленные данные",
		})
		return
	}
	log.Println("New registry try:", req.Login, req.Password)

	password, err := service.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось хэшировать пароль",
		})
		return
	}
	message, err := db.Register(req.Login, password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось создать пользователя",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
