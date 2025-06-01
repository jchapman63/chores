package main

import (
	"context"
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
	l := log.New(os.Stdout, "CHORES: ", log.LstdFlags)
	l.Println("Starting Chores Application...")
	cronLog := cron.VerbosePrintfLogger(l)

	// Load database configuration
	dbConfig := config.LoadDBConfig()
	connString := dbConfig.GetDBConnectionString()

	// Establish database connection
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		l.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		l.Fatalf("Unable to ping database: %v", err)
	}
	l.Println("Successfully connected to database")

	// Initialize sqlc queries
	queries := db.New(dbPool)

	// Initialize rotation service with database queries
	rotationService := rotation.NewService(queries)
	l.Println("Rotation service initialized successfully")

	// Create a new cron scheduler with seconds field disabled
	c := cron.New(cron.WithLogger(cronLog))
	// Add a job that runs every Monday at 9am to rotate chores automatically
	_, err = c.AddFunc("0 9 * * 1", func() {
		// Create context for the rotation operation
		ctx := context.Background()

		cronLog.Info("Running scheduled chore rotation...")
		// Rotate the chores using the rotation service
		if err := rotationService.RotateChores(ctx); err != nil {
			cronLog.Info("Error rotating chores: %v", err)
		} else {
			cronLog.Info("Chores rotated successfully")
		}
	})
	if err != nil {
		l.Println(err, "Error setting up chore rotation cron job")
	}

	// Start a test cron that runs every minute
	_, err = c.AddFunc("* * * * *", func() {
		cronLog.Info("Test cron job running every minute")
	})
	if err != nil {
		cronLog.Error(err, "Error setting up test cron job")
	}

	// Start the cron scheduler in the background
	c.Start()
	l.Println("Chore rotation scheduler started")
	l.Println("Application running. Press Ctrl+C to exit.")

	// Set up signal handling for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	s := <-sig
	l.Printf("Received signal %s. Shutting down...\n", s)
	c.Stop()
	l.Println("Cron scheduler stopped. Goodbye!")
}
