package main

import (
	"html/template"
	"strconv"
	"vkr/internal/handler"
	"vkr/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const port = 8085

const sessionSecret = "Xqg1j1tmGupISrzeTbzTb6DsvZLscOZ5"
const sessionLifetime = 60 * 60 * 1
const sessionName = "userSession"

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{MaxAge: sessionLifetime})

	router.Use(sessions.Sessions(sessionName, store))

	router.Static("/web", "web")

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/basic/*.html"))
	router.SetHTMLTemplate(tmpl)

	router.GET("/", handler.IndexHandler)
	router.GET("/login", handler.GetLoginHandler)
	router.POST("/login", handler.PostLoginHandler)
	router.GET("/register", handler.GetRegisterHandler)
	router.POST("/register", handler.PostRegisterHandler)
	auth := router.Group("/")
	auth.Use(service.AuthRequired())
	{
		auth.POST("/logout", handler.PostLogoutHandler)
		auth.GET("/profile", handler.GetProfileHandler)
	}
	admin := router.Group("/")
	admin.Use(service.AdminRequired())
	{
		admin.GET("/admin", handler.GetAdminHandler)
	}

	router.Run(":" + strconv.Itoa(port))
}
