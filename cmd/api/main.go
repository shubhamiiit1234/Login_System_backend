package main

import (
	"context"
	"fmt"
	"log"
	"login/internal/app"
	"login/internal/database"
	"login/internal/router"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Login System!")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Note: No .env file found, using system environment variables")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}

	connection := fmt.Sprintf("host=%s port=5432 user=postgres password=mysecretpassword dbname=Login sslmode=disable", dbHost)
	err = database.InitializeDB(connection)
	if err != nil {
		fmt.Println("error initializing database: ", err.Error())
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		fmt.Println("error getting database instance: ", err.Error())
		return
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	rdb := database.InitializeRedis(redisAddr)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis failed to connect: %v", err)
	}
	defer rdb.Close()

	handlersContainer := app.NewHandlersContainer(db, rdb)

	r := router.NewRouter(handlersContainer)
	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("error starting server: ", err.Error())
		return
	}
}
