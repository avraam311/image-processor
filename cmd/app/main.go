package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	handlers "github.com/avraam311/image-processor/internal/api/handlers/images"
	"github.com/avraam311/image-processor/internal/api/server"
	repository "github.com/avraam311/image-processor/internal/repository/images"
	service "github.com/avraam311/image-processor/internal/service/images"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"

	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go"
)

const (
	configFilePath = "config/local.yaml"
	envFilePath    = ".env"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zlog.Init()
	val := validator.New()
	cfg := config.New()
	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load env file")
	}
	cfg.EnableEnv("")
	if err := cfg.LoadConfigFiles(configFilePath); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config file")
	}

	opts := &dbpg.Options{
		MaxOpenConns:    cfg.GetInt("db.max_open_conns"),
		MaxIdleConns:    cfg.GetInt("db.max_idle_conns"),
		ConnMaxLifetime: cfg.GetDuration("db.conn_max_lifetime"),
	}
	slavesDNSs := []string{}
	masterDNS := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.GetString("DB_USER"), cfg.GetString("DB_PASSWORD"),
		cfg.GetString("DB_HOST"), cfg.GetString("DB_PORT"),
		cfg.GetString("DB_NAME"), cfg.GetString("DB_SSL_MODE"),
	)
	db, err := dbpg.New(masterDNS, slavesDNSs, opts)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	kafkaProd := kafka.NewProducer(cfg.GetStringSlice("kafka.brokers"), cfg.GetString("kafka.topic"))
	minioEndpoint := cfg.GetString("MINIO_HOST" + ":" + "MINIO_PORT")
	minioAccessKey := cfg.GetString("MINIO_ACCESS_KEY")
	minioSecret := cfg.GetString("MINIO_SECRET")
	minioSSL := cfg.GetBool("MINIO_SSL")
	minioClient, err := minio.New(minioEndpoint, minioAccessKey, minioSecret, minioSSL)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to initialize minio client")
	}
	repo := repository.NewRepository(db)
	srvc := service.NewService(repo, kafkaProd, cfg, minioClient)
	hand := handlers.NewHandler(srvc, val)

	router := server.NewRouter(cfg.GetString("api.gin_mode"), hand)
	srv := server.NewServer(cfg.GetString("server.port"), router)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to run server")
		}
	}()
	zlog.Logger.Info().Msg("server is running")

	<-ctx.Done()
	zlog.Logger.Info().Msg("shutdown signal received")

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	zlog.Logger.Info().Msg("shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to shutdown server")
	}
	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		zlog.Logger.Info().Msg("timeout exceeded, forcing shutdown")
	}

	if err := db.Master.Close(); err != nil {
		zlog.Logger.Printf("failed to close master DB: %v", err)
	}
	for i, s := range db.Slaves {
		if err := s.Close(); err != nil {
			zlog.Logger.Printf("failed to close slave DB %d: %v", i, err)
		}
	}

	if err := kafkaProd.Close(); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to close kafka producer")
	}
}
