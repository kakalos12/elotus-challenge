package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// Helper function to extract and validate token
func getAuthenticatedUser(c *gin.Context) (*Claims, error) {
	tokenString, err := c.Cookie("token")
	if err != nil || tokenString == "" {
		return nil, errors.New("no token provided")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Explicitly check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := getAuthenticatedUser(c)
		if err != nil {
			// On any auth error (expired, invalid, missing), clear cookie and redirect
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func redirectIfAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := getAuthenticatedUser(c)
		if err == nil && claims != nil {
			// User is authenticated
			if c.Request.Method == http.MethodGet {
				c.Redirect(http.StatusFound, "/upload")
			} else {
				c.String(http.StatusForbidden, "Forbidden: You are already logged in")
			}
			c.Abort()
			return
		}
		// Not authenticated, proceed
		c.Next()
	}
}
