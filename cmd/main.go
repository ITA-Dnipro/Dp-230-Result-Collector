package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/app"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/config"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/mongodb"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/service"
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(run(cfg))
}

func run(cfg config.Config) error {
	app, err := app.NewApp(cfg)
	if err != nil {
		return err
	}

	validate := validator.New()
	reportRepo := mongodb.NewReportMongoRepo(app.MongoClient.Client)
	service := service.NewReportService(reportRepo, validate)
	pb.RegisterReportServiceServer(app.Server.Server, service)

	return app.Run(context.Background())
}
