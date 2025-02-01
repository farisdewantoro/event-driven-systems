package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"loanservice/configs"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
)

// Worker processes the Asynq task by calling `/send-email`
func HandleEmailVerificationTask(ctx context.Context, t *asynq.Task) error {
	var payload map[string]string
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Call API to send email
	apiURL := "http://localhost:8080/send-email"
	reqBody, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call send-email API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send-email API returned status: %d", resp.StatusCode)
	}

	log.Printf("Email verification request sent for UserID: %s", payload["user_id"])
	return nil
}

func NewWorker() {
	cfg := configs.Get()

	// Asynq Worker Setup
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Address},
		asynq.Config{Concurrency: 5},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:verify", HandleEmailVerificationTask)

	log.Println("Worker started, waiting for tasks...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("Could not start worker: %v", err)
	}
}
