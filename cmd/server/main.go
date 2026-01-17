package main

import (
	"aiksava-lb/internal/router"
	"aiksava-lb/internal/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
	
	services.InitBackends()
	services.ColdStart()
	services.StartSessionCleanup()

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Load Balancer running on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}