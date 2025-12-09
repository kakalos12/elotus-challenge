package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	// Retry connection logic for docker-compose startup
	for range 10 {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database after retries:", err)
	}

	createTables()
	log.Println("Connected to database successfully")
}

func createTables() {
	userTable := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	fileTable := `CREATE TABLE IF NOT EXISTS file_uploads (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		filename VARCHAR(255),
		content_type VARCHAR(255),
		size BIGINT,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	_, err = db.Exec(fileTable)
	if err != nil {
		log.Fatal("Error creating file_uploads table:", err)
	}
}
