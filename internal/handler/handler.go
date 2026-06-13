package handler

import (
	"log"
	"net/http"
	"vkr/internal/db"
	"vkr/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var req Credentials
var value Val

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Val struct {
	Login   string
	IsLogin bool
	IsAdmin bool
}

func setValue(c *gin.Context) {
	session := sessions.Default(c)
	login, ok := session.Get("login").(string)
	is_login, ok := session.Get("is_login").(bool)
	is_admin, ok := session.Get("is_admin").(bool)
	if !ok {
		log.Println("Значения не установлены")
	}

	value = Val{
		Login:   login,
		IsLogin: is_login,
		IsAdmin: is_admin,
	}
}

func IndexHandler(c *gin.Context) {
	setValue(c)
	c.HTML(http.StatusOK, "main.html", gin.H{
		"login":    value.Login,
		"is_login": value.IsLogin,
		"is_admin": value.IsAdmin,
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

		service.SetSession(c, user)

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login successful",
			"redirect": "/",
		})
	} else {
		log.Println("Не верный пароль")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
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
		"message":  message,
		"redirect": "/login",
	})
}

func GetAdminHandler(c *gin.Context) {
	setValue(c)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"login":    value.Login,
		"is_login": value.IsLogin,
		"is_admin": value.IsAdmin,
	})
}

func PostLogoutHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("New logout try")

	service.UnsetSession(c)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Logout successful",
		"redirect": "/",
	})
}

func GetProfileHandler(c *gin.Context) {
	setValue(c)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"login":    value.Login,
		"is_login": value.IsLogin,
		"is_admin": value.IsAdmin,
	})
}
