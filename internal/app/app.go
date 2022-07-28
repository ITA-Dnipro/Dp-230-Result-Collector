package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ITA-Dnipro/Dp-230-Result-Collector/internal/config"
	grpcWrapper "github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/grpc"
	"github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/kafka"
	mongoWrapper "github.com/ITA-Dnipro/Dp-230-Result-Collector/pkg/mongodb"
)

type Component interface {
	Start(context.Context) error
	Stop(context.Context) error
	Name() string
}

type App struct {
	config     config.Config
	logger     *zap.Logger
	components []Component

	//Components
	Server      *grpcWrapper.ServerWrapper
	MongoClient *mongoWrapper.ClientWrapper
	Producer    *kafka.SyncProducerWrapper
}

func NewApp(cfg config.Config, logger *zap.Logger) (*App, error) {

	app := &App{
		logger: logger,
		config: cfg,
	}
	if err := app.initComponents(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initComponents() error {
	a.components = make([]Component, 0)

	srv, err := grpcWrapper.NewServer("grpc server", a.config.Server, a.logger)
	if err != nil {
		return err
	}
	a.Server, a.components = srv, append(a.components, srv)

	client, err := mongoWrapper.NewClient("mongo client", a.config.MongoDB, a.logger)
	if err != nil {
		return err
	}
	a.MongoClient, a.components = client, append(a.components, client)

	producer, err := kafka.NewSyncProducer("producer", a.config.Producer, a.logger)
	if err != nil {
		return err
	}
	a.Producer, a.components = producer, append(a.components, producer)

	return a.validateComponents()
}

func (a *App) validateComponents() error {
	v := reflect.ValueOf(a).Elem()
	t := reflect.TypeOf(a).Elem()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)
		if ft.IsExported() && fv.IsNil() {
			return errors.New(fmt.Sprintf("%s not initialized", ft.Name))
		}
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()
		signaled := make(chan os.Signal, 1)
		signal.Notify(signaled, os.Interrupt, syscall.SIGTERM)
		select {
		case s := <-signaled:
			a.logger.Info("exiting", zap.Stringer("signal", s))
			signal.Stop(signaled)

		case <-ctx.Done():
			signal.Stop(signaled)
		}
	}()

	for _, c := range a.components {
		if err := c.Start(ctx); err != nil {
			a.logger.Warn("failed to start", zap.String("component", c.Name()), zap.Error(err))
			return fmt.Errorf("can't start component [%s]: %w", c.Name(), err)
		}
	}

	<-ctx.Done()
	if err := a.stop(); err != nil {
		a.logger.Warn("failed to stop app", zap.Error(err))
	}
	return ctx.Err()
}

func (a *App) stop() error {
	shutdownTimeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	for _, c := range a.components {
		err := c.Stop(ctx)
		if err != nil {
			a.logger.Warn("failed to stop", zap.String("component", c.Name()), zap.Error(err))
		}
	}
	return nil
}

func verifyConfigStructOrPtr(cfg interface{}) error {
	cfgVal := reflect.ValueOf(cfg)
	if cfgVal.Kind() == reflect.Ptr {
		if cfgVal.IsNil() {
			return errors.New("application configuration should be non-nil")
		}
		cfgVal = cfgVal.Elem()
	}
	if cfgVal.Kind() != reflect.Struct {
		return errors.New("application configuration should be a struct")
	}
	return nil
}
