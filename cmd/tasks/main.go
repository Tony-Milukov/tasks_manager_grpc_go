package main

import (
	"fmt"
	"log/slog"
	"os"
	"sso_3.0/internal/app"
	configParser "sso_3.0/internal/config"
)

func main() {
	fmt.Println("Hello World!")
	config := configParser.MustGetConfig()

	fmt.Printf("enviroment: %s", config.Env)
	log := getLogger()
	//
	////Setup APp
	app, err := app.New(config, log)
	//
	if err != nil {
		panic(fmt.Errorf("errors: %e", err))
	}

	//start grpc server
	app.GrpcServer.MustRun()
}

func getLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
