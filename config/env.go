package config

import (
	"github.com/joho/godotenv"
	"os"
	"log"
)

func GetEnv(key string) string {
	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }