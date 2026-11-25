// cmd/main.go is wiring only — no business logic. See
// document/ARCHITECTURE.md "Request flow". Pattern taken from
// ekyc-service's cmd/main.go (the most recently written of the 8 services).
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"example.com/TODO-service-name/config"
	"example.com/TODO-service-name/internal/handler/health"
	"example.com/TODO-service-name/internal/middleware"
	"example.com/TODO-service-name/internal/routes"
	"example.com/TODO-service-name/internal/server"
	pkglogger "example.com/TODO-service-name/pkg/logger"
)

func main() {
	cfg := initConfig()

	if err := pkglogger.Initialize(pkglogger.Config{
		Level:       cfg.Logger.Level,
		Encoding:    cfg.Logger.Encoding,
		ServiceName: cfg.Server.ServiceID,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "logger init failed: %v\n", err)
		os.Exit(1)
	}
	defer pkglogger.Sync() //nolint:errcheck
	logger := pkglogger.Get()

	// Construct real-vs-stub adapters here based on config/env, once —
	// never branch on "am I in a test" inside business logic. e.g.:
	//   storage := storage.New(cfg.Storage, logger)   // --storage s3
	//   cache := cache.New(cfg.CacheDatabase, logger)  // --cache redis

	router := newRouter(cfg)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := server.NewServer(addr, router, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)

	serveWithGracefulShutdown(srv, addr, logger)
}

func newRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		middleware.RequestIDs(),
		gin.Recovery(),
	)
	router.Use(func(c *gin.Context) {
		c.Set("service_id", cfg.Server.ServiceID)
		c.Next()
	})

	healthHandler := health.NewHandler(cfg.Server.ServiceID)
	routes.RegisterHealthRoutes(router, healthHandler)
	// routes.RegisterExampleRoutes(router, hExample)  // TODO: wire your domain routes

	return router
}

// serveWithGracefulShutdown runs the HTTP server until SIGINT/SIGTERM, then
// drains in-flight requests before exiting.
func serveWithGracefulShutdown(srv *http.Server, addr string, logger *zap.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting server", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
	wg.Wait()
	logger.Info("server stopped")
}

func initConfig() *config.Config {
	cfg, err := config.Load(os.Getenv("CONFIG_FILE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}
	return cfg
}
