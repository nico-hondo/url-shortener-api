package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/nico-hondo/url-shortener/CMD/docs"
	"github.com/nico-hondo/url-shortener/internal/config"
	"github.com/nico-hondo/url-shortener/internal/router"
)

// @title           Url Shortener
// @version         1.0
// @description     Completed Eventhub backend during Koda Bootcamp
// @host            localhost:8000
// @BasePath        /
// @securityDefinitions.apikey  BearerToken
// @in                          header
// @name                        Authorization
// @description                 Bearer Token used as identity for accessing backend
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

	// connect to redis
	rdb := config.NewRdbConfig(os.Getenv("RDB_USER"), os.Getenv("RDB_PASS"), os.Getenv("RDB_ADDR"), os.Getenv("RDB_PORT"))
	redisClient := rdb.ConnectRdb()

	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Println("Redis is not ready\nReason: ", err.Error())
		return
	}
	log.Println("Redis Ready")

	// Generate gin Engine
	r := gin.Default()

	// Deklarasi Router (endpoint & method HTTP)
	router.InitMainRouter(r, pool, redisClient)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
