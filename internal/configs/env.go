package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	PORT       string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBSchema   string
	JwtSecret  string
	ServerENV  string
}

var Env *env

func SetupENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err.Error())
	}

	Env = &env{
		PORT:       os.Getenv("PORT"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBSchema:   os.Getenv("DB_SCHEMA"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
		ServerENV:  os.Getenv("SERVER_ENV"),
	}
}
