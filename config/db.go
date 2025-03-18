package config

import "os"

type Database struct {
	ConnURL string
}

func LoadDatabaseConfig() Database {
	return Database{
		ConnURL: os.Getenv("DB_CONN_URL"),
	}
}
