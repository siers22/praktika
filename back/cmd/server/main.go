package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/siers22/praktika/back/internal/config"
	"github.com/siers22/praktika/back/internal/db"
	"github.com/siers22/praktika/back/internal/handler"
	"github.com/siers22/praktika/back/internal/middleware"
	"github.com/siers22/praktika/back/internal/repository"
	"github.com/siers22/praktika/back/internal/router"
	"github.com/siers22/praktika/back/internal/service"
)

func main() {
	cfg := config.Load()

	// Structured JSON logging
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	slog.Info("starting inventory API",
		"port", cfg.ServerPort,
		"debug", cfg.Debug,
	)

	// Database
	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Migrations
	if err := db.RunMigrations(pool, "migrations"); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	// Seed default admin if no users exist
	if err := db.SeedAdminUser(pool, cfg.DefaultAdminUser, cfg.DefaultAdminPass); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		slog.Error("failed to create upload dir", "error", err)
		os.Exit(1)
	}

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	equipmentRepo := repository.NewEquipmentRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	departmentRepo := repository.NewDepartmentRepository(pool)
	inventoryRepo := repository.NewInventoryRepository(pool)
	movementRepo := repository.NewMovementRepository(pool)
	auditRepo := repository.NewAuditRepository(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, auditRepo, cfg)
	userSvc := service.NewUserService(userRepo, auditRepo)
	equipmentSvc := service.NewEquipmentService(equipmentRepo, auditRepo, cfg.UploadDir)
	categorySvc := service.NewCategoryService(categoryRepo, auditRepo)
	departmentSvc := service.NewDepartmentService(departmentRepo, auditRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo, equipmentRepo, auditRepo)
	movementSvc := service.NewMovementService(movementRepo, equipmentRepo, auditRepo)
	reportSvc := service.NewReportService(equipmentRepo, movementRepo, inventoryRepo)
	auditSvc := service.NewAuditService(auditRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	equipmentHandler := handler.NewEquipmentHandler(equipmentSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	departmentHandler := handler.NewDepartmentHandler(departmentSvc)
	inventoryHandler := handler.NewInventoryHandler(inventorySvc)
	movementHandler := handler.NewMovementHandler(movementSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	auditHandler := handler.NewAuditHandler(auditSvc)

	// Middleware
	authMw := middleware.NewAuthMiddleware(authSvc)

	// Router
	r := router.Setup(
		authHandler,
		userHandler,
		equipmentHandler,
		categoryHandler,
		departmentHandler,
		inventoryHandler,
		movementHandler,
		reportHandler,
		auditHandler,
		authMw,
		cfg.UploadDir,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	slog.Info("server stopped")
}
