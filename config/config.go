package config

type Config struct {
	Database   Database
	Jwt        Jwt
	Server     Server
	App        App
	Pagination Pagination
}

type App struct {
	Env      string
	LogLevel string
}

func LoadConfig() Config {
	return Config{
		Database:   LoadDatabaseConfig(),
		Jwt:        LoadJwtConfig(),
		Server:     LoadServerConfig(),
		App:        LoadAppConfig(),
		Pagination: LoadPaginationConfig(),
	}
}
