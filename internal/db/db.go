package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Подключение базы данных
func DB() (*sqlx.DB, error) {
	//Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_SSLMODE"))

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	return db, err
}
