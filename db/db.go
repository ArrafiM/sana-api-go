package db

import (
	"fmt"
	"sana-api/config"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var CON *gorm.DB

func ConnecDatabase() {
	host := config.GetEnv("DB_HOST")
	user := config.GetEnv("DB_USER")
	pass := config.GetEnv("DB_PASS")
	name := config.GetEnv("DB_NAME")
	port, _ := strconv.Atoi(config.GetEnv("DB_PORT"))
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		host, user, pass, name, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	CON = database
}
