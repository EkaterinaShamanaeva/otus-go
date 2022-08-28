package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	configuration "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/init_storage"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	internalhttp "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := configuration.NewConfig()

	if err := config.BuildConfig(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatalf("Logger error: %v", err)
	}

	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Database.Username,
		config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name,
		config.Database.SSLMode)

	storage, err := init_storage.NewStorage(ctx, config.Storage, dsn)
	if err != nil {
		logg.Error("failed to connect DB: " + err.Error())
	}
	defer storage.Close(ctx)

	logg.Info("DB connected...")

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)
	grpc := internalgrpc.New(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err = server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err = grpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}

		if err = calendar.Close(ctx); err != nil {
			logg.Error("failed close storage: " + err.Error())
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(config.Server.Host, config.Server.HttpPort)
		if err = server.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(config.Server.Host, config.Server.GrpcPort)
		if err = grpc.Start(ctx, addrServer); err != nil {
			logg.Error("failed to start gRPC server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
