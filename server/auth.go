package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

func registerPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func registerHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.String(http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error hashing password")
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		c.String(http.StatusConflict, "Error creating user (username might be taken)")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func loginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.String(http.StatusBadRequest, "Invalid request payload")
		return
	}

	var storedPassword string
	var userID int
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", user.Username).Scan(&userID, &storedPassword)
	if err == sql.ErrNoRows {
		c.String(http.StatusUnauthorized, "Invalid username or password")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		c.String(http.StatusUnauthorized, "Invalid username or password")
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating token")
		return
	}

	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	c.String(http.StatusOK, fmt.Sprintf("Logged in successfully! Token: %s", tokenString))
}
