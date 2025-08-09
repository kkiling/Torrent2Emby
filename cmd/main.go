package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/kkiling/goplatform/config"
	platformserver "github.com/kkiling/goplatform/server"

	appconfig "github.com/kkiling/torrent2emby/internal/config"
	"github.com/kkiling/torrent2emby/internal/container"
	"github.com/kkiling/torrent2emby/internal/server"
)

func main() {
	var args config.Arguments
	if _, err := flags.Parse(&args); err != nil {
		log.Fatal(err)
	}

	cfgProvider, err := config.NewProvider(args)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := appconfig.NewEnvConfig(cfgProvider)
	if err != nil {
		log.Fatal(err)
	}

	cn, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = cn.MkvMergePipeline().StartMergePipeline(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err = cn.GetContentDelivery().Complete(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	logger := cn.GetLogger()

	srv := server.NewTorrent2EmbyServer(
		logger,
		platformserver.Config{
			Host:                    cfg.Server.Host,
			GrpcPort:                cfg.Server.GrpcPort,
			HttpPort:                cfg.Server.HttpPort,
			MaxSendMessageLength:    cfg.Server.MaxSendMessageLength,
			MaxReceiveMessageLength: cfg.Server.MaxReceiveMessageLength,
			ShutdownTimeout:         cfg.Server.ShutdownTimeout,
		},
		cn.GetTvShowLibrary(),
		cn.GetContentDelivery(),
	)
	go func() {
		err = srv.Start(ctx)
		if err != nil {
			log.Fatalf("fail start app: %v", err)
		}
	}()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Infof("--- shutdown application ---")
		cancel()
	}()

	<-ctx.Done()
	logger.Infof("--- stopped application ---")
	srv.Stop()
	logger.Infof("--- stop application ---")
}
