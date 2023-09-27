package main

import (
	auth_service "auth-service"
	"auth-service/pkg/handler"
	"auth-service/pkg/repository"
	"auth-service/pkg/service"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	ctx := context.Background()
	client, err := repository.NewMongoDB(
		ctx, repository.Config{
			Port: viper.GetString("dbport"),
			Host: viper.GetString("host"),
		},
	)

	if err != nil {
		logrus.Fatalf("failed initialize client: %s", err.Error())
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logrus.Fatalf("failed uninitialize client: %s", err.Error())
		}
	}()

	repos := repository.NewRepository(ctx, client)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(auth_service.Server)
	if err := srv.Run(viper.GetString("port"), http.HandlerFunc(handlers.Serve)); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
