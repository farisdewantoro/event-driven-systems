package cmd

import (
	"eventdrivensystem/configs"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "event-driven system",
	Short: "event-driven system",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func init() {
	configs.Load()
	rootCmd.AddCommand(apiServerCmd)
	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(outboxWorkerCmd)
	rootCmd.AddCommand(asynqWorkerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
