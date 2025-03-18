package config

import "os"

type Server struct {
	Port string
}

func LoadServerConfig() Server {
	return Server{
		Port: os.Getenv("SERVER_PORT"),
	}
}
