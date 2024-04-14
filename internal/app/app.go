package app

import (
	"avito-test-task/internal/app/routes"
	"avito-test-task/internal/config"
	"avito-test-task/internal/lib/jwt"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/service"
	"avito-test-task/internal/storage/pgs"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App is the main application structure. It holds all the dependencies and the server.
type App struct {
	logger        *slog.Logger
	jwtManager    *jwt.Manager
	bannerService *service.Banner
}

// New creates a new instance of the App.
func New(logger *slog.Logger, jwtManager *jwt.Manager, bannerSvc *service.Banner) *App {
	return &App{
		logger:        logger,
		jwtManager:    jwtManager,
		bannerService: bannerSvc,
	}
}

// startServer starts the handlers server.
func (a *App) startServer(ctx context.Context, server *http.Server) {
	a.logger.Info("starting server", slog.String("address", server.Addr))
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("listen and serve returned err", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	a.shutdownGracefully(ctx, server)
}

// shutdownGracefully shuts down the server gracefully.
func (a *App) shutdownGracefully(ctx context.Context, server *http.Server) {
	a.logger.Info("gracefully shutting down")
	waitForReturn(
		ctx,
		10*time.Second,
		server.Shutdown,
		func() { a.logger.Error("failed to shutdown server") },
	)
	a.logger.Info("server stopped")
}

// Run starts the application with the default parameters.
func Run() {
	cfg := config.MustLoad(os.Args[1:], os.LookupEnv)
	logger := initLogger(cfg.Env)
	storage := initStorage(context.Background(), cfg.DB.ConnectionString(), logger)
	jwtManager := jwt.NewManager(string(cfg.JwtSettings.SecretKey), time.Duration(cfg.JwtSettings.Expire))
	bannerService := service.NewBannerService(storage, storage, storage, storage, logger)
	app := New(logger, jwtManager, bannerService)
	run(context.Background(), cfg, app)
	waitForReturn(
		context.Background(),
		10*time.Second,
		storage.Close,
		func() { logger.Error("failed to close storage") },
	)
}

// RunWithConfig starts the application with the provided configuration.
// It is used mainly for testing.
func RunWithConfig(
	ctx context.Context,
	args []string,
	getenv func(string) (string, bool),
	app *App,
) {
	cfg := config.MustLoad(args, getenv)
	run(ctx, cfg, app)
}

// run starts the app.
func run(ctx context.Context, cfg *config.Config, app *App) {
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      routes.New(app.logger, app.jwtManager, app.bannerService),
		WriteTimeout: time.Duration(cfg.HTTPServer.Timeout),
		IdleTimeout:  time.Duration(cfg.HTTPServer.IdleTimeout),
		ReadTimeout:  time.Duration(cfg.HTTPServer.Timeout),
	}

	app.startServer(ctx, server)
}

// waitForReturn waits for the provided function to return, but only for the provided duration.
func waitForReturn(
	ctx context.Context,
	duration time.Duration,
	waitFunc func(ctx context.Context) error,
	timeoutCallback func(),
) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	done := make(chan struct{})
	go func() {
		_ = waitFunc(ctx)
		done <- struct{}{}
		close(done)
	}()

	select {
	case <-ctx.Done():
		timeoutCallback()
	case <-done:
		return
	}
}

// initLogger initializes the logger based on the environment.
func initLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.LocalEnv:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.ProdEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

// initStorage initializes the application storage.
func initStorage(ctx context.Context, connString string, logger *slog.Logger) *pgs.Storage {
	storage, err := pgs.New(ctx, connString)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("storage initialized", slog.String("storage", "postgres"))
	return storage
}
