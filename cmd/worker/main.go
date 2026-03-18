package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"path/filepath"
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"
	"project-mini-e-commerce/pkg/mail"
	"project-mini-e-commerce/pkg/rabbitmq"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMQ    rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	cfg         *config.Config
	logger      *zerolog.Logger
}

func NewWorker(cfg *config.Config) *Worker {
	log := utils.NewLoggerWithPath("rabbitmq.log", "info")
	rabbitMQ, err := rabbitmq.NewRabbitMQService(utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"), log)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to create RabbitMQ service")
	}

	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to create mail provider factory")
		return nil
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to create mail service")
		return nil
	}

	return &Worker{
		rabbitMQ:    rabbitMQ,
		mailService: mailService,
		cfg:         cfg,
		logger:      log,
	}
}

func (w *Worker) StartWorker(ctx context.Context) error {
	const emailQueue = "email_queue"
	handler := func(body []byte) error {
		w.logger.Debug().Msgf("Received message: %s", string(body))
		var email mail.Email

		if err := json.Unmarshal(body, &email); err != nil {
			w.logger.Error().Err(err).Msg("Failed to unmarshal email message")
			return err
		}

		if err := w.mailService.SendMail(ctx, &email); err != nil {
			return utils.WrapError(err, "Failed to send reset password email", utils.ErrCodeInternal)
		}

		w.logger.Info().Msgf("Email sent to %s with subject: %s", email.To, email.Subject)
		return nil
	}
	if err := w.rabbitMQ.Consume(ctx, emailQueue, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to consume messages from RabbitMQ")
		return err
	}

	w.logger.Info().Msgf("Worker started and consuming messages from RabbitMQ %s", emailQueue)

	<-ctx.Done()
	w.logger.Info().Msg("worker stopping due to context cancellation: ")
	return ctx.Err()
}

func (w *Worker) StopWorker(ctx context.Context) error {
	w.logger.Info().Msg("worker is shutting down .....")
	if err := w.rabbitMQ.Close(); err != nil {
		w.logger.Error().Err(err).Msg("Failed to close RabbitMQ connection")
		return err
	}
	w.logger.Info().Msg("Worker stopped and RabbitMQ connection closed")

	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			w.logger.Warn().Msg("worker shutdown timed out, forcing exit")
			return ctx.Err()
		}
	default:
	}

	w.logger.Info().Msg("Worker shutdown complete")
	return nil
}

func main() {
	rootDir := utils.GetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.Config{
		Level:       "info",
		Filename:    logFile,
		MaxSize:     1,
		MaxAge:      5,
		MaxBackups:  5,
		Compress:    true,
		Environment: common.Environment(utils.GetEnv("APP_ENV", "development")),
	})

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Logger.Warn().Msg("No .env file found")
	} else {
		logger.Logger.Info().Msg(".env file loaded successfully in worker")
	}
	configFile := config.NewConfig()

	worker := NewWorker(configFile)
	if worker == nil {
		logger.Logger.Fatal().Msg("Failed to initialize worker")
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := worker.StartWorker(ctx); err != nil && err != context.Canceled {
			logger.Logger.Fatal().Err(err).Msg("Worker failed to start")
		}
	}()

	<-ctx.Done()
	logger.Logger.Info().Msg("Worker is shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := worker.StopWorker(shutdownCtx); err != nil {
		logger.Logger.Error().Err(err).Msg("Worker shutdown encountered an error")
	}
	wg.Wait()
	logger.Logger.Info().Msg("Worker shutdown complete")
}
