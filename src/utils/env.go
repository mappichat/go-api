package utils

import "os"

type EnvironmentVariables struct {
	PORT    string
	DB_HOST string
}

var Env *EnvironmentVariables = &EnvironmentVariables{}

func ConfigureEnv() {
	if Env.PORT = os.Getenv("PORT"); Env.PORT == "" {
		Env.PORT = "8080"
	}
	if Env.DB_HOST = os.Getenv("DB_HOST"); Env.DB_HOST == "" {
		Env.DB_HOST = ""
	}
}
