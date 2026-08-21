package serverconfig

import (
	"github.com/joho/godotenv"
)

type Config struct{
	ServerPort string
	DatabaseURL string
	Environment string
}

func LoadConfig() (*Config, error){
	if err := godotenv.Load()
}
