package main

import (
	"fmt"
	"login/internal/app"
	"login/internal/database"
	"login/internal/router"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Login System!")

	connection := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=Login sslmode=disable"
	err := database.InitializeDB(connection)
	if err != nil {
		fmt.Println("error initializing database: ", err.Error())
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		fmt.Println("error getting database instance: ", err.Error())
		return
	}

	handlersContainer := app.NewHandlersContainer(db)

	r := router.NewRouter(handlersContainer)
	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("error starting server: ", err.Error())
		return
	}
}
