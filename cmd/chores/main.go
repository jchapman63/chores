package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	"github.com/jchapman63/chores/config"
	db "github.com/jchapman63/chores/internal/db/sqlc"
	"github.com/jchapman63/chores/internal/rotation"
)

func main() {
	fmt.Println("Starting Chores Application...")

	// Load database configuration
	dbConfig := config.LoadDBConfig()
	connString := dbConfig.GetDBConnectionString()

	// Establish database connection
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	fmt.Println("Successfully connected to database")

	// Initialize sqlc queries
	queries := db.New(dbPool)

	// Initialize rotation service with database queries
	rotationService := rotation.NewService(queries)
	fmt.Println("Rotation service initialized successfully")

	// Create a new cron scheduler with seconds field disabled
	c := cron.New(cron.WithSeconds())
	// Add a job that runs every Monday at 9am to rotate chores automatically
	_, err = c.AddFunc("0 9 * * 1", func() {
		// Create context for the rotation operation
		ctx := context.Background()
		
		log.Println("Running scheduled chore rotation...")
		// Rotate the chores using the rotation service
		if err := rotationService.RotateChores(ctx); err != nil {
			log.Printf("Error rotating chores: %v", err)
		} else {
			log.Println("Chores rotated successfully")
		}
	})
	if err != nil {
		log.Fatalf("Error setting up chore rotation cron job: %v", err)
	}

	// Start the cron scheduler in the background
	c.Start()
	fmt.Println("Chore rotation scheduler started")

	fmt.Println("Application running. Press Ctrl+C to exit.")
	
	// Set up signal handling for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	s := <-sig
	fmt.Printf("Received signal %s. Shutting down...\n", s)
	c.Stop()
	fmt.Println("Cron scheduler stopped. Goodbye!")
}
