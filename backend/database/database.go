package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var logger = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Llongfile)
var loggerError = log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Llongfile)
var tableCreated bool = true

func init() {
	err := godotenv.Load()
	if err != nil {
		loggerError.Print(err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connectionString := "postgresql://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	for i := 0; i < 5; i++ {
		DB, err = sql.Open("postgres", connectionString)
		if err != nil {
			loggerError.Printf("Error connecting to database (attempt %d): %v", i+1, err)
			time.Sleep(5 * time.Second)
		} else {
			logger.Printf("Connected to database successfully")
			break
		}
	}

	if DB == nil {
		loggerError.Fatal("Failed to connect to database after multiple attempts")
	}

	createTable()
}

func createTable() {
	if tableCreated {
		logger.Println("Table 'shortened_urls' already exists")
		return
	}

	logger.Println("Creating table 'shortened_urls' if not exists")
	sqlStatement := `
        CREATE TABLE IF NOT EXISTS shortened_urls (
            id SERIAL PRIMARY KEY,
            key VARCHAR(6) UNIQUE NOT NULL,
            long_url TEXT
        );
    `

	_, err := DB.Exec(sqlStatement)
	if err != nil {
		loggerError.Print("Error creating table:", err)
	} else {
		logger.Println("Table 'shortened_urls' created successfully")
		tableCreated = true
	}
}
