package main

import (
	"log"
	"os"

	"github.com/durgaprasad97005/GoFiberAssignment2/src"
)

func main() {
	app := src.SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server is running at port:" + port)

	log.Fatal(app.Listen(":" + port))
}