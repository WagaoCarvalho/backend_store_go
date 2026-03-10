package config

import "os"

type Server struct {
	Port    string
	BaseURL string
	IDPath  string
}

func LoadServerConfig() Server {
	return Server{
		Port:    os.Getenv("SERVER_PORT"),
		BaseURL: os.Getenv("API_BASE_URL"),
		IDPath:  os.Getenv("API_ID_PATH"),
	}
}
