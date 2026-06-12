package main

import (
	"strconv"
	"vkr/internal/handler"

	"github.com/gin-gonic/gin"
)

const Port = 8080

func main() {
	router := gin.Default()
	router.Static("/web", "web")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handler.IndexHandler)
	router.GET("/login", handler.GetLoginHandler)
	//router.POST("/login", handler.PostLoginHandler)
	router.GET("/register", handler.RegisterHandler)

	router.Run(":" + strconv.Itoa(Port))
}
