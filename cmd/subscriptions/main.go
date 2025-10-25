package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	config2 "github.com/s-turchinskiy/EffectiveMobile/cmd/subscriptions/config"
	"github.com/s-turchinskiy/EffectiveMobile/internal/common/closer"
	"github.com/s-turchinskiy/EffectiveMobile/internal/handlers"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/repository"
	"github.com/s-turchinskiy/EffectiveMobile/internal/repository/postgresql"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

var outputPathsLogs string

func init() {

	err := config2.InitializePublicConfig()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&outputPathsLogs, "logs", config2.PublicConfig.OutputPathsLogsDefault, "Paths to output (file, stdout) logs")
	flag.Parse()

	if err := logger.Initialize(outputPathsLogs); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(outputPathsLogs)

}
func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var closer = &closer.Closer{}

	errorsCh := make(chan error)
	defer close(errorsCh)

	doneCh := make(chan struct{})

	err := godotenv.Load(config2.PublicConfig.EnvFilename)
	if err != nil {
		logger.Log.Debugw("Error loading .env file", "error", err.Error())
	}

	configData, err := config2.GetConfig()
	if err != nil {
		logger.Log.Errorw("Get Settings error", "error", err.Error())
		errorsCh <- err
	}

	rep, err := postgresql.NewPostgresStorage(ctx, configData.Database.String(), configData.Database.DBName)
	if err != nil {
		logger.Log.Debugw("Connect to database error", "error", err.Error())
		log.Fatal(err)
	}

	//rep = repository.NewRepositoryWithRetry(rep, config.PublicConfig.HTTPServer.RetryStrategy)
	closer.Add(rep.Close)

	httpServer := getHTTPServer(ctx, rep, configData.Address.String())
	closer.Add(funcHTTPServerShutdown(httpServer))

	go runHTTPServer(httpServer, configData.Address.String(), errorsCh)

	<-ctx.Done()

	err = shutdown(closer, doneCh)
	if err != nil {
		logger.Log.Fatalw("fatal error", "error", err.Error())
	}

}

func getHTTPServer(ctx context.Context, repository repository.Repository, addr string) *http.Server {

	handler := handlers.NewHandler(
		ctx,
		repository,
		config2.PublicConfig.HTTPServer.RetryStrategy)

	router := handlers.Router(handler)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  config2.PublicConfig.HTTPServer.ReadTimeout,
		WriteTimeout: config2.PublicConfig.HTTPServer.WriteTimeout,
	}

	return server

}

func runHTTPServer(server *http.Server, addr string, errorsCh chan error) {

	logger.Log.Info("Running server", zap.String("address", addr))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {

		logger.Log.Errorw("Server startup error", "error", err.Error())
		errorsCh <- err
		return
	}

}

func funcHTTPServerShutdown(httpServer *http.Server) func(ctx context.Context) error {

	return func(ctx context.Context) error {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.Log.Infow("HTTP server stopped with error", zap.String("error", err.Error()))
		} else {
			logger.Log.Infow("HTTP server stopped")
		}
		return err
	}
}
func shutdown(closer *closer.Closer, doneCh chan struct{}) error {

	log.Println("shutting down server gracefully")

	close(doneCh)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config2.PublicConfig.ShutdownTimeout)
	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer: %v", err)
	}

	return nil

}
