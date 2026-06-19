package main

import (
	"html/template"
	"log"
	"strconv"
	"vkr/internal/config"
	"vkr/internal/handler"
	"vkr/internal/repository"
	"vkr/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
}

func main() {
	defer repository.CloseDB()

	router := gin.Default()

	store := cookie.NewStore([]byte(config.SessionSecret))
	lifetime, err := strconv.Atoi(config.SessionLifetime)
	if err != nil {
		log.Println("Warning: session lifetime not a number")
	}
	store.Options(sessions.Options{MaxAge: lifetime})

	router.Use(sessions.Sessions(config.SessionName, store))

	router.Static("/web", "web")
	router.Static("/uploads", "./uploads")

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/basic/*.html"))
	router.SetHTMLTemplate(tmpl)

	router.GET("/", handler.IndexHandler)
	router.GET("/login", handler.GetLoginHandler)
	router.GET("/register", handler.GetRegisterHandler)
	router.POST("/login", handler.PostLoginHandler)
	router.POST("/register", handler.PostRegisterHandler)

	auth := router.Group("/")
	auth.Use(service.AuthRequired())
	{
		auth.GET("/profile", handler.GetProfileHandler)
		auth.GET("/images", handler.GetImagesHandler)
		auth.POST("/logout", handler.PostLogoutHandler)
		auth.POST("/upload", handler.PostUploadHandler)
	}

	admin := router.Group("/")
	admin.Use(service.AdminRequired())
	{
		admin.GET("/admin", handler.GetAdminHandler)
	}

	router.Run(":" + config.Port)
}
