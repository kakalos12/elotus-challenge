package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	initDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Public routes (redirect if already logged in)
	public := r.Group("/")
	public.Use(redirectIfAuthenticatedMiddleware())
	{
		public.GET("/register", registerPage)
		public.POST("/register", registerHandler)

		public.GET("/login", loginPage)
		public.POST("/login", loginHandler)
	}

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
