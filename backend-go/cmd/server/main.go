// Command server is the entry point for the Rentacar CRM Go backend
// (Gin + PostgreSQL + Redis).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/instaagrammeta/rentacar/backend-go/internal/auth"
	"github.com/instaagrammeta/rentacar/backend-go/internal/config"
	"github.com/instaagrammeta/rentacar/backend-go/internal/database"
	"github.com/instaagrammeta/rentacar/backend-go/internal/handler"
	"github.com/instaagrammeta/rentacar/backend-go/internal/media"
	"github.com/instaagrammeta/rentacar/backend-go/internal/middleware"
	"github.com/instaagrammeta/rentacar/backend-go/internal/scheduler"
	"github.com/instaagrammeta/rentacar/backend-go/internal/seed"
	"github.com/instaagrammeta/rentacar/backend-go/internal/service"
)

func main() {
	cfg := config.Load()
	if err := cfg.EnsureDirectories(); err != nil {
		log.Fatalf("failed to prepare directories: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration: %v", err)
	}

	rds, err := database.ConnectRedis(cfg)
	if err != nil {
		log.Printf("warning: Redis unavailable, continuing without cache: %v", err)
		rds = nil
	}

	// Seed default data on first launch (admin + settings, demo data optional).
	if envBool("RENTACAR_AUTO_SEED", true) {
		seed.Run(db, envBool("RENTACAR_SEED_DEMO", true))
	}

	tokens := auth.NewManager(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	gen := media.NewGenerator(cfg.UploadDir, cfg.ExportDir)
	svc := service.New(db, cfg, rds, tokens, gen)
	authMw := &middleware.AuthService{Tokens: tokens, Redis: rds}
	h := handler.New(svc, authMw, cfg)

	cronJobs := scheduler.Start(svc)
	if cronJobs != nil {
		defer cronJobs.Stop()
	}

	router := h.SetupRouter()
	srv := &http.Server{
		Addr:    cfg.HTTPHost + ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("Rentacar CRM backend listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
}

func envBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "yes")
}
