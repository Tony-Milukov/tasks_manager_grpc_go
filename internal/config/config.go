package configParser

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env         string
	DbUrl       string
	GrpcPort    string
	TokenSecret string
}

func MustGetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUrl := getEnv("DB_URL")
	fmt.Printf("%s", dbUrl)
	env := getEnv("ENV")
	grpcPort := getEnv("GRPC_PORT")
	tokenSecret := getEnv("GRPC_PORT")

	return &Config{
		env,
		dbUrl,
		grpcPort,
		tokenSecret,
	}

}

func getEnv(key string) string {
	env := os.Getenv(key)

	if env == "" {
		panic(fmt.Sprintf("the env %s was not set", key))
	}

	return env
}
