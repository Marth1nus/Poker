package main

import (
	"fmt"
	"go-poker/routes"
)

func main() {
	fmt.Println("Starting server http://localhost:8080")
	routes.StartServer(8080)
}
