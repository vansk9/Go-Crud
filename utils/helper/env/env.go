package env

import "os"

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

func IsLocal() bool {
	return os.Getenv("ENV") == "local"
}
