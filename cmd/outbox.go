package cmd

import (
	"context"
	"fmt"
	"loanservice/configs"
	"loanservice/internal/models"
	"loanservice/pkg/databases"
	"loanservice/pkg/util"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-openapi/strfmt"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var outboxWorkerCmd = &cobra.Command{
	Use:   "outbox-worker",
	Short: "Runs the outbox worker",
	Run: func(cmd *cobra.Command, args []string) {
		outboxWorker := NewOutBoxWorker()

		// Create a cancelable context for graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // Make sure to cancel the context when we're done

		wg := &sync.WaitGroup{}
		go outboxWorker.Run(ctx, wg)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Received shutdown signal, initiating graceful shutdown...")
		cancel() // Cancel the context to initiate a graceful shutdown

		// Wait for all workers to finish
		wg.Wait()
		log.Println("Graceful shutdown complete.")
	},
}

type OutboxWorker struct {
	db         *gorm.DB
	cfg        *configs.AppConfig
	queue      *asynq.Client
	workerPool chan struct{}
}

// ProcessOutboxJobs processes jobs with limited concurrency
func (o *OutboxWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down outbox worker...")
			return
		default:
			o.processOutboxJobs(wg)
		}
	}
}

func (o *OutboxWorker) processOutboxJobs(wg *sync.WaitGroup) error {
	tx := o.db.Begin()

	delayNextIteration := time.Duration(rand.Intn(o.cfg.Outbox.DurationIntervalInMs)) * time.Millisecond

	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		time.Sleep(delayNextIteration)
		return tx.Error
	}

	// Fetch 100 messages to process using raw SQL
	var outboxes []models.Outbox
	err := tx.Raw(`
			SELECT id, status, attempt, execute_at
			FROM outbox
			WHERE status IN (?, ?) AND execute_at <= ?
			ORDER BY execute_at asc
			FOR UPDATE SKIP LOCKED
			LIMIT ?
		`, models.OutboxStatusPending, models.OutboxStatusRetrying, time.Now(), o.cfg.Outbox.MaxBatchSize).Scan(&outboxes).Error

	if err != nil {
		tx.Rollback()
		log.Printf("Error fetching data from outbox: %v", err)
		time.Sleep(delayNextIteration)
		return err
	}

	if len(outboxes) == 0 {
		log.Printf("No outboxes to process")
		tx.Rollback()
		// Sleep for the jitter time
		fmt.Printf("Sleeping for %v...\n", delayNextIteration)
		time.Sleep(delayNextIteration)
		return nil
	}

	log.Printf("Found %d outboxes to process", len(outboxes))

	// Update all selected outboxes to PROCESSING status in the same transaction
	err = o.setStatusProcessing(tx, outboxes)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating outboxes to PROCESSING: %v", err)
		time.Sleep(10 * time.Second)
		return err
	}

	// Commit the transaction before processing
	err = tx.Commit().Error
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		time.Sleep(delayNextIteration)
		return err
	}

	// Process the outboxes concurrently with a worker pool
	for _, outbox := range outboxes {
		wg.Add(1)
		o.workerPool <- struct{}{} // Acquire a worker slot
		go func(outbox models.Outbox) {
			defer func() {
				<-o.workerPool
				wg.Done()
			}() // Release worker slot after processing

			bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			o.processMessage(bgCtx, outbox)
		}(outbox)
	}

	fmt.Printf("Sleeping for %v...\n", delayNextIteration)
	time.Sleep(delayNextIteration)

	return nil
}

// processMessage processes the message, retries on failure, and updates its status
func (o *OutboxWorker) processMessage(ctx context.Context, outbox models.Outbox) {
	var err error

	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Error processing message %s: %v", outbox.ID, r)
			return
		}
		if err != nil {
			tx.Rollback()
			log.Printf("Error processing message %s: %v", outbox.ID, err)
			return
		}
		err = tx.Commit().Error
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}()
	// Simulate message processing, can be replaced with actual processing logic (e.g., sending to a queue)
	errProcess := o.process(ctx, outbox)
	finishedProcessTime := time.Now()

	if errProcess != nil {
		outbox.ErrorMessage = util.ToPointer(errProcess.Error())
		// If processing fails, update status to RETRYING
		if outbox.Attempt >= int64(o.cfg.Outbox.MaxRetries) {
			outbox.Status = models.OutboxStatusFailed
			// If max retries are reached, update status to FAILED
			err = tx.Updates(outbox).Error
			if err != nil {
				log.Printf("Error updating status for message %s: %v", outbox.ID, err)
				return
			}
			return
		}
		// Compute the next retry time using Exponential Backoff
		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = 5 * time.Minute // Start with 1 min delay
		expBackoff.MaxInterval = 10 * time.Minute    // Max delay between retries
		expBackoff.MaxElapsedTime = 1 * time.Hour    // Stop retrying after 1 hour
		expBackoff.Multiplier = 2.0                  // Exponential growth
		expBackoff.RandomizationFactor = 0.5         // Add jitter to prevent sync issues
		expBackoff.Reset()                           // Reset to ensure a fresh calculation

		// Calculate new backoff time based on attempt count
		for i := 0; i < int(outbox.Attempt+1); i++ {
			expBackoff.NextBackOff() // Advance to correct attempt time
		}
		nextExecuteAt := time.Now().Add(expBackoff.NextBackOff())

		outbox.ExecuteAt = nextExecuteAt
		outbox.Status = models.OutboxStatusRetrying

		err = tx.Updates(outbox).Error

		if err != nil {
			log.Printf("Error updating status for message %s: %v", outbox.ID, err)
		}
		return
	}

	// Successfully processed, update status to SENT
	outbox.Status = models.OutboxStatusSent
	outbox.SentAt = &finishedProcessTime
	err = tx.Updates(outbox).Error
	if err != nil {
		log.Printf("Error updating status for message %s: %v", outbox.ID, err)
		return
	}

}

// process simulates message processing, replace this with your actual processing logic
func (o *OutboxWorker) process(ctx context.Context, outbox models.Outbox) error {
	// Create a random number to simulate processing logic
	randomNumber := rand.Intn(10) + 1 // Generate a number between 1 and 10
	if randomNumber == 1 {
		return fmt.Errorf("random processing error")
	}
	return nil
}

func (o *OutboxWorker) setStatusProcessing(tx *gorm.DB, outboxes []models.Outbox) error {
	outboxIds := make([]strfmt.UUID4, len(outboxes))
	for i, outbox := range outboxes {
		outboxIds[i] = outbox.ID
	}

	qUpdate := "UPDATE outbox SET status = ?, attempt = attempt + 1 WHERE id IN ?"
	err := tx.Exec(qUpdate, models.OutboxStatusProcessing, outboxIds).Error

	if err != nil {
		log.Printf("error set status processing for outbox_ids: %v err: %v", outboxIds, err)
		return err
	}
	return nil
}

// NewOutBoxWorker initializes and returns the OutboxWorker
func NewOutBoxWorker() *OutboxWorker {
	cfg := configs.Get()
	db, err := databases.NewSqlDb(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Redis.Address})

	// Initialize the worker pool with a defined concurrency limit
	workerPool := make(chan struct{}, cfg.Outbox.MaxConcurrency) // Limit concurrency to 10 workers

	return &OutboxWorker{
		db:         db,
		cfg:        cfg,
		queue:      asynqClient,
		workerPool: workerPool,
	}
}
