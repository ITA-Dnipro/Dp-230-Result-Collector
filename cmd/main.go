package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/app"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/config"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/kafka"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/mongodb"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/service"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/usecase"
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
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	app, err := app.NewApp(cfg, logger)
	if err != nil {
		return err
	}

	validate := validator.New()
	repository := mongodb.NewReportMongoRepo(app.MongoClient.Client)
	producer := kafka.NewReportProducer("test", app.Producer.SyncProducer, logger)
	usecases := usecase.NewReportUsecase(repository, validate, producer)
	service := service.NewReportService(usecases)
	pb.RegisterReportServiceServer(app.Server.Server, service)

	return app.Run(context.Background())
}
