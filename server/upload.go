package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 8 << 20 // 8 MB

func uploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func uploadHandler(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized: User not found in context")
		return
	}
	userClaims := claims.(*Claims)
	userID := userClaims.UserID

	err := c.Request.ParseMultipartForm(maxUploadSize)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to parse form: %v", err))
		return
	}

	file, header, err := c.Request.FormFile("data")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get 'data' file from form")
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		c.String(http.StatusBadRequest, "File is too large (max 8MB)")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !isImageContentType(contentType) {
		c.String(http.StatusBadRequest, "Invalid content type. Only images are allowed.")
		return
	}

	tempDir := "/tmp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			log.Printf("Failed to create temp directory %s: %v", tempDir, err)
			c.String(http.StatusInternalServerError, "Failed to create temporary directory")
			return
		}
	}

	tempFileName := fmt.Sprintf("%d-%s-%s", userID, time.Now().Format("20060102150405"), header.Filename)
	tempFilePath := filepath.Join(tempDir, tempFileName)

	out, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Failed to create temporary file %s: %v", tempFilePath, err)
		c.String(http.StatusInternalServerError, "Failed to create temporary file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Failed to save file to %s: %v", tempFilePath, err)
		c.String(http.StatusInternalServerError, "Failed to save uploaded file")
		return
	}

	stmt, err := db.Prepare("INSERT INTO file_uploads (user_id, filename, content_type, size, uploaded_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare DB statement: %v", err)
		c.String(http.StatusInternalServerError, "Failed to prepare database for metadata storage")
		return

	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, header.Filename, contentType, header.Size, time.Now())
	if err != nil {
		log.Printf("Failed to insert file metadata into DB: %v", err)
		c.String(http.StatusInternalServerError, "Failed to store file metadata")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File uploaded successfully. Filename: %s, Size: %d, Content-Type: %s", tempFileName, header.Size, contentType))
}

func isImageContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}
