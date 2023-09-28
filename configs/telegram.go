package configs

import (
	"os"
)

// err := godotenv.Load()
// if err != nil {
// 	log.Fatal("Error loading .env file")
// }
var TELEGRAM_FOLLOOOW_TOKEN = os.Getenv("TELEGRAM_FOLLOOOW_TOKEN")
var TELEGRAM_FOLLOOOW_CHANNEL = os.Getenv("TELEGRAM_FOLLOOOW_CHANNEL")

// func EnvTelegram() (string, string) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	return os.Getenv("MONGO_URI"), os.Getenv("")
// }

// func EnvMongoDB() string {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	return os.Getenv("MONGO_DB")
// }
