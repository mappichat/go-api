package utils

import "os"

type Env struct {
	PORT string
}

var env *Env

func configureEnv() {
	env.PORT = os.Getenv("PORT")
	if env.PORT == "" {
		env.PORT = "8080"
	}
}
