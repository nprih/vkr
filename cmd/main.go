package main

import (
	"vkr/internal/handler"

	"github.com/gin-gonic/gin"
)

//func main() {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/", handler.IndexHandler)
//	mux.HandleFunc("GET /login", handler.GetLoginHandler)
//	mux.HandleFunc("POST /login", handler.PostLoginHandler)
//	mux.HandleFunc("/register", handler.RegisterHandler)
//
//	log.Println("Starting server...")
//	if err := http.ListenAndServe(":8080", mux); err != nil {
//		log.Println(err)
//	}
//}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handler.IndexHandler)
	router.GET("/login", handler.GetLoginHandler)
	//router.POST("/login", handler.PostLoginHandler)
	router.GET("/register", handler.RegisterHandler)

	router.Run(":8080")
}
