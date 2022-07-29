package config

import (
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/grpc"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/kafka"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/mongodb"
)

// Config takes config parameters from environment, or uses default.
type Config struct {
	Server   *grpc.ServerConfiguration
	MongoDB  *mongodb.ClientConfiguration
	Producer *kafka.ProducerConfiguration
}
