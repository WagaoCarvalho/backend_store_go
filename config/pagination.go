package config

import (
	"os"
	"strconv"
)

type Pagination struct {
	DefaultLimit  int
	DefaultOffset int
	MaxLimit      int
}

func LoadPaginationConfig() Pagination {
	return Pagination{
		DefaultLimit:  getEnvAsInt("PAGINATION_DEFAULT_LIMIT", 50),
		DefaultOffset: getEnvAsInt("PAGINATION_DEFAULT_OFFSET", 0),
		MaxLimit:      getEnvAsInt("PAGINATION_MAX_LIMIT", 500),
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
