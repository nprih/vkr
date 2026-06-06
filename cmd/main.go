package main

import (
	"vkr/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handler.IndexHandler)
	router.GET("/login", handler.GetLoginHandler)
	//router.POST("/login", handler.PostLoginHandler)
	router.GET("/register", handler.RegisterHandler)

	router.Run(":8080")
}
