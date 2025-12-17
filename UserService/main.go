package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"user_central/server"
	"user_central/services"
	"user_central/storage"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed loading from .env, %s\n", err.Error())
		return
	}
	db, err := setupDB()
	if err != nil {
		log.Printf("failed connecting to database, %s\n", err.Error())
		return
	}
	repo := storage.PostgresUserRepo{Database: db}
	err = repo.Init()
	if err != nil {
		log.Printf("failed initializing repository, %s\n", err.Error())
		return
	}

	service := &services.GinUserService{Repo: &repo}

	server := server.ServerInit(service)

	if err = server.Run(os.Getenv("PORT")); err != nil {
		log.Printf("failed to run application %s\n", err.Error())
		return
	}

}

func setupDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("warning: couldn't load .env file: %v", err)
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with database")
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}
