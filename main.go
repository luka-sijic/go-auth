package main

import (
	"app/database"
	"app/server"
)

func main() {
	database.Connect()
	defer database.Close()

	server.Start()
}