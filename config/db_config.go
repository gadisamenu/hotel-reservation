package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB_NAME string

var MONGO_DB_URI string
var JWT_SECRET string
var HTTP_LISTEN_ADDRESS string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal(err)
		}
	}

	DB_NAME = os.Getenv("MONGO_DB_NAME")
	MONGO_DB_URI = os.Getenv("MONGO_DB_URI")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	HTTP_LISTEN_ADDRESS = os.Getenv("HTTP_LISTEN_ADDRESS")
}
