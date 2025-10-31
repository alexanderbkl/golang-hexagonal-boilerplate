package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	pb "github.com/alexanderbkl/golang-hexagonal-boilerplate/api/grpc"
	dbadapter "github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/adapters/db"
	gqladapter "github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/adapters/graphql"
	grpcadapter "github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/adapters/grpc"
	redisadapter "github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/adapters/redis"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/internal/services"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/pkg/config"
	"github.com/alexanderbkl/golang-hexagonal-boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	grpcServer "google.golang.org/grpc"
)

func main() {
	log := logger.New()
	log.Info("Starting application...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	log.Info("Connecting to database...")
	dbPool, err := pgxpool.New(context.Background(), cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Info("Database connection established")

	// Initialize Redis connection
	log.Info("Connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Info("Redis connection established")

	// Initialize repositories
	userRepo := dbadapter.NewPostgresRepository(dbPool)
	cacheRepo := redisadapter.NewRedisRepository(redisClient)

	// Initialize services
	userService := services.NewUserService(userRepo, cacheRepo)

	// Start gRPC server
	go func() {
		log.Infof("Starting gRPC server on port %s...", cfg.Server.GRPCPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcSrv := grpcServer.NewServer()
		pb.RegisterUserServiceServer(grpcSrv, grpcadapter.NewUserServiceServer(userService))

		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start GraphQL/HTTP server
	log.Infof("Starting GraphQL server on port %s...", cfg.Server.HTTPPort)

	resolver := gqladapter.NewResolver(userService)
	srv := handler.NewDefaultServer(gqladapter.NewExecutableSchema(gqladapter.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.HTTPPort),
		Handler: http.DefaultServeMux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Infof("GraphQL playground available at http://localhost:%s/", cfg.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
