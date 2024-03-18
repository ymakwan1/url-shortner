package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB
var logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
var loggerError = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
var tableCreated bool // Flag to track if table creation has been attempted

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

	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		loggerError.Print("Error connecting to database:", err)
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
