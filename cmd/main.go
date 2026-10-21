package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nico-hondo/url-shortener/internal/config"
	"github.com/nico-hondo/url-shortener/internal/router"
)

func main() {
	// load env seawal mungkin
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
		return
	}

	// connect ke DB
	pdb := config.NewPsqlDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready\nReason: ", err.Error())
		return
	}

	log.Println("Database Ready")

	// Generate gin Engine
	r := gin.Default()

	// Deklarasi Router (endpoint & method HTTP)
	router.InitMainRouter(r, pool)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
