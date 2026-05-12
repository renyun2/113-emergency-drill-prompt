package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"emergency-drill/internal/database"
	"emergency-drill/internal/handlers"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if strings.EqualFold(os.Getenv("GIN_MODE"), "release") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	handlers.Register(r, db)
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	log.Println("API " + listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
