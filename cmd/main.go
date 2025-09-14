package main

import (
	"log"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/configs"
	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/routers"
)

// @title Tickitz Booking API
// @version 1.0
// @description API for book ticket cinemas
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	db, err := configs.InitDB()
	if err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	rdb, err := configs.InitRedis()
	if err != nil {
		log.Println("RDB init failed:", err)
	}

	if rdb != nil {
		defer rdb.Close()
	}

	r := routers.MainRouter(db, rdb)

	r.Run(":8080")
}
