package utils

import "os"

type EnvironmentVariables struct {
	PORT       string
	DB_HOST    string
	JWT_SECRET string
}

var Env *EnvironmentVariables = &EnvironmentVariables{}

func ConfigureEnv() {
	if Env.PORT = os.Getenv("PORT"); Env.PORT == "" {
		Env.PORT = "8080"
	}
	if Env.DB_HOST = os.Getenv("DB_HOST"); Env.DB_HOST == "" {
		Env.DB_HOST = "localhost:9042,localhost:9043,localhost:9044"
	}
	if Env.JWT_SECRET = os.Getenv("JWT_SECRET"); Env.JWT_SECRET == "" {
		Env.JWT_SECRET = "secret"
	}
}
