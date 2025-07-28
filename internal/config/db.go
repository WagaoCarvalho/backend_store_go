package config

import (
	"os"
)

type Database struct {
	ConnURL string
}

var LoadDatabaseConfig = func() Database {

	return Database{
		ConnURL: os.Getenv("DB_CONN_URL"),
	}
}
