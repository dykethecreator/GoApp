package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dykethecreator/GoApp/internal/auth/handler"
	"github.com/dykethecreator/GoApp/internal/auth/middleware"
	"github.com/dykethecreator/GoApp/internal/auth/service"
	"github.com/dykethecreator/GoApp/internal/auth/store"
	"github.com/dykethecreator/GoApp/pkg/database"
	appjwt "github.com/dykethecreator/GoApp/pkg/jwt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Environment loading strategy:
	// 1) If APP_ENV is set, try to load .env.{APP_ENV}
	// 2) If RUNNING_IN_DOCKER is set, default to .env.docker
	// 3) Otherwise, default to .env.local
	// 4) Finally, load base .env if present (for overrides)
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		if os.Getenv("RUNNING_IN_DOCKER") != "" {
			appEnv = "docker"
		} else {
			appEnv = "local"
		}
	}

	// Try environment-specific file first
	_ = godotenv.Load(".env." + appEnv)
	// Then load base .env optionally to allow local overrides
	_ = godotenv.Load()

	// Veritabanı bağlantısını kur
	db, err := database.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Determine gRPC port (default 50051)
	grpcPort := os.Getenv("AUTH_SERVICE_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	listenAddr := fmt.Sprintf(":%s", grpcPort)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Build a TokenManager for interceptor (same secret/durations as service)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	tm, err := appjwt.NewTokenManager(jwtSecret, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		log.Fatalf("failed to init token manager: %v", err)
	}

	// Create gRPC server with auth interceptor; AuthService methods are exempt inside the interceptor.
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor(tm)),
	)

	// Bağımlılıkları oluştur (DI - Dependency Injection)
	userStore := store.NewUserStore(db.DB)
	deviceStore := store.NewUserDeviceStore(db.DB)
	authService := service.NewAuthService(userStore, deviceStore)
	authHandler := handler.NewAuthHandler(authService)

	// Handler'ı gRPC sunucusuna kaydet
	authHandler.Register(s)

	log.Printf("auth_service listening on %s (env=%s)", listenAddr, os.Getenv("APP_ENV"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
