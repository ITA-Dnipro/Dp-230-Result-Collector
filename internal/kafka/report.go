package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/model"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type reportProducer struct {
	topic    string
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewReportProducer(topic string, producer sarama.SyncProducer, log *zap.Logger) *reportProducer {
	return &reportProducer{
		topic:    topic,
		producer: producer,
		logger:   log,
	}
}

func (p *reportProducer) Send(r *model.Report) error {
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("producer marshal report: %w", err)
	}
	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(b),
	})
	p.logger.Info("successfully sent message:",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
	return err
}
