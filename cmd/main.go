package cmd

import (
	"fmt"
	"loanservice/configs"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loan-service",
	Short: "loan-service is a service to manage loans",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func init() {
	configs.Load()
	rootCmd.AddCommand(apiServerCmd)
	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(outboxWorkerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
