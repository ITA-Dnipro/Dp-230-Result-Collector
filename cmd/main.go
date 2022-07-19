package main

import (
	"context"
	"log"
	"net"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/config"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/mongodb"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/service"
	pb "github.com/ITA-Dnipro/Dp-230-Result-Collector/proto"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("xss", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mongoDBConn, err := mongodb.NewMongoDB(ctx, cfg)
	if err != nil {
		log.Fatal("cannot connect mongodb", err)
	}
	defer func() {
		if err := mongoDBConn.Disconnect(ctx); err != nil {
			log.Fatal("mongoDBConn.Disconnect", err)
		}
	}()

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reportRepo := mongodb.NewReportMongoRepo(mongoDBConn)
	service := service.NewReportService(reportRepo)

	srv := grpc.NewServer()
	pb.RegisterReportServiceServer(srv, service)

	log.Printf("starting server %s", cfg.Server.Port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
