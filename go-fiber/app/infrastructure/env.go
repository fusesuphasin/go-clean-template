package infrastructure

import (
	"github.com/fusesuphasin/go-fiber/app/interfaces"
	"github.com/joho/godotenv"
)

// Load is load configs from a env file.
func Load(logger interfaces.Logger) {
	err := godotenv.Load()
	if err != nil {
		logger.LogError("%s", err)
	}
}