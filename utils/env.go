package utils

import "os"

type EnvironmentVariables struct {
	PORT string
}

var Env *EnvironmentVariables = &EnvironmentVariables{}

func ConfigureEnv() {
	Env.PORT = os.Getenv("PORT")
	if Env.PORT == "" {
		Env.PORT = "8080"
	}
}
