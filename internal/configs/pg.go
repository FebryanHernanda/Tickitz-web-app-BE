package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env\nCause: ", err.Error())
	}

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASSWORD")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")
	connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println("Connection string:", connstring)

	db, err := pgxpool.New(context.Background(), connstring)
	if err != nil {
		log.Println("Failed to connect to database\nCause: ", err.Error())
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Println("Ping to DB failed\nCause: ", err.Error())
	} else {
		log.Println("DB Connected")
	}

	return db, nil
}
