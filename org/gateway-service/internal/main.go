package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Elbujito/2112/src/graphql-api/go/graph"
	"github.com/go-redis/redis/v8"
	"github.com/rs/cors"
)

var redisClient *redis.Client

func main() {
	// Environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	serverPort := getEnv("SERVER_PORT", "4000")

	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})
	defer redisClient.Close()

	// Initialize custom resolver
	resolver := NewCustomResolver(redisClient)

	// GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// HTTP server with CORS
	mux := http.NewServeMux()
	mux.Handle("/query", srv)
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	handlerWithCORS := cors.Default().Handler(mux)

	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: handlerWithCORS,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 Server ready at http://localhost:%s", serverPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped.")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
