package utils

import (
	"os"
	"strconv"
)

type EnvironmentVariables struct {
	PORT                     string
	DB_CONNECTION_STRING     string
	JWT_SECRET               string
	MAX_RESOLUTION           int
	VOTE_DISTANCE_MULTIPLIER float64
	AUTH0_EXTENSION_SECRET   string
}

var Env *EnvironmentVariables = &EnvironmentVariables{}

func ConfigureEnv() {
	var err error
	if Env.PORT = os.Getenv("PORT"); Env.PORT == "" {
		Env.PORT = "8080"
	}
	if Env.DB_CONNECTION_STRING = os.Getenv("DB_CONNECTION_STRING"); Env.DB_CONNECTION_STRING == "" {
		Env.DB_CONNECTION_STRING = "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
	}
	if Env.JWT_SECRET = os.Getenv("JWT_SECRET"); Env.JWT_SECRET == "" {
		Env.JWT_SECRET = "secret"
	}
	if Env.MAX_RESOLUTION, err = strconv.Atoi(os.Getenv("MAX_RESOLUTION")); err != nil {
		Env.MAX_RESOLUTION = 6
	}
	if Env.VOTE_DISTANCE_MULTIPLIER, err = strconv.ParseFloat(os.Getenv("VOTE_DISTANCE_MULTIPLIER"), 64); err != nil {
		Env.VOTE_DISTANCE_MULTIPLIER = 0.5
	}
	if Env.AUTH0_EXTENSION_SECRET = os.Getenv("AUTH0_EXTENSION_SECRET"); Env.AUTH0_EXTENSION_SECRET == "" {
		Env.AUTH0_EXTENSION_SECRET = "secret"
	}
}
