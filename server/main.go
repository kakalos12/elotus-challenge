package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	initDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/register", registerPage)
	r.POST("/register", registerHandler)

	r.GET("/login", loginPage)
	r.POST("/login", loginHandler)

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
