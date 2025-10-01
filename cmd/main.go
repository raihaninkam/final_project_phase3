package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/raihaninkam/finalPhase3/internals/configs"
	"github.com/raihaninkam/finalPhase3/internals/routers"
)

// @title 					Social Media
// @version 				1.0
// @host						localhost:3009
// @basePath				/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan format: Bearer <token>
func main() {
	// Inisialization databae for this project
	db, err := configs.InitPg()
	if err != nil {
		log.Println("FAILED TO CONNECT DB")
		return
	}

	defer db.Close()

	err = configs.PingDB(db)
	if err != nil {
		log.Println("PING TO DB FAILED", err.Error())
		return
	}

	log.Println("DB CONNECTED")

	// inisialization redis
	rdb := configs.InitRedis()
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		log.Println("failed ping on redis \nCause:", cmd.Err().Error())
		return
	}
	log.Println("Redis Connected")
	defer rdb.Close()

	// Inisialization engine gin, HTTP framework
	router := routers.InitRouter(db, rdb)
	router.Run(":3009")
}
