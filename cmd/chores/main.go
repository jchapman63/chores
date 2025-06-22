package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"

	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jchapman63/chores/config"
	db "github.com/jchapman63/chores/internal/db/sqlc"
	"github.com/jchapman63/chores/internal/rotation"
	"github.com/jchapman63/chores/internal/sns"
)

type Onboardee struct {
	Name  string
	Email string
	Chore string
}

func main() {
	l := log.New(os.Stdout, "CHORES: ", log.LstdFlags)
	l.Println("Starting Chores Application...")
	cronLog := cron.VerbosePrintfLogger(l)

	ctx := context.Background()

	// Load configuration
	cfg := config.LoadConfig()

	// Establish SNS client
	snsClient, err := sns.NewSNSClient(ctx)
	if err != nil {
		l.Fatalf("Failed to create SNS client: %v", err)
	}

	// Create a new database connection pool
	dbPool, err := pgxpool.New(context.Background(), cfg.DB.GetDBConnectionString())
	if err != nil {
		l.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		l.Fatalf("Unable to ping database: %v", err)
	}

	// Initialize sqlc queries and rotation service
	queries := db.New(dbPool)
	rotationService := rotation.NewService(queries)
	if err := rotationService.InitializeChores(ctx); err != nil {
		l.Fatalf("Failed to initialize chores: %v", err)
	}

	// send the initial chore digest to SNS
	rms, err := rotationService.GetRoommates(ctx)
	if err != nil {
		l.Fatalf("Failed to get roommates: %v", err)
	}
	_, err = snsClient.Client.Publish(ctx, &awssns.PublishInput{
		Message:  rotationService.CreateChoreDigest(rms),
		TopicArn: &cfg.AWS.SNSTopicARN,
	})
	if err != nil {
		l.Printf("Failed to get publish initial digest: %v", err)
	}

	// Create a new cron scheduler with seconds field disabled
	c := cron.New(cron.WithLogger(cronLog))
	// Add a job that runs every Monday at 9am to rotate chores automatically
	_, err = c.AddFunc("0 9 * * 1", func() {
		cronLog.Info("Running scheduled chore rotation...")
		// Rotate the chores using the rotation service
		rms, err := rotationService.RotateChores(ctx)
		if err != nil {
			cronLog.Info("Error rotating chores: %v", err)
		}
		_, err = snsClient.Client.Publish(ctx, &awssns.PublishInput{
			Message:  rotationService.CreateChoreDigest(rms),
			TopicArn: &cfg.AWS.SNSTopicARN,
		})
		if err != nil {
			cronLog.Info("Error sending chore digest: %v", err)
		}
	})
	if err != nil {
		l.Println(err, "Error setting up chore rotation cron job")
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
