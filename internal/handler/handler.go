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

func setValue(c *gin.Context, user *db.User) {
	session := sessions.Default(c)
	session.Set("login", user.Login)
	session.Set("is_login", true)
	session.Set("is_admin", user.Is_admin)
	session.Save()

	login, ok := session.Get("login").(string)
	if !ok {
		log.Println("Значение login не установлено")
	}
	is_login, ok := session.Get("is_login").(bool)
	if !ok {
		log.Println("Значение is_login не установлено")
	}
	is_admin, ok := session.Get("is_admin").(bool)
	if !ok {
		log.Println("Значение is_admin не установлено")
	}
	value = Val{
		Login:   login,
		IsLogin: is_login,
		IsAdmin: is_admin,
	}
}

func unsetValue(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("login")
	session.Delete("is_login")
	session.Delete("is_admin")
	session.Save()
	value = Val{}
}

func IndexHandler(c *gin.Context) {
	//session := sessions.Default(c)

	//login, ok := session.Get("login").(string)
	//if !ok {
	//	log.Println("Значение login не установлено")
	//}
	//is_login, ok := session.Get("is_login").(bool)
	//if !ok {
	//	log.Println("Значение is_login не установлено")
	//}
	//is_admin, ok := session.Get("is_admin").(bool)
	//if !ok {
	//	log.Println("Значение is_admin не установлено")
	//}

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
		setValue(c, user)
		//session := sessions.Default(c)
		//session.Set("login", user.Login)
		//session.Set("is_login", true)
		//session.Set("is_admin", user.Is_admin)
		//session.Save()

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
	c.HTML(http.StatusOK, "main.html", gin.H{
		"title": "Main website",
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

	unsetValue(c)
	//session := sessions.Default(c)
	//session.Delete("login")
	//session.Delete("is_login")
	//session.Delete("is_admin")
	//session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message":  "Logout successful",
		"redirect": "/",
	})
}
