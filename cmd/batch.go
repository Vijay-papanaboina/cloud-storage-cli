/*
Copyright © 2025 vijay papanaboina

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Batch job management commands",
	Long: `Manage and monitor batch jobs.

Available commands:
  status - Get batch job status and progress`,
}

// batchStatusCmd represents the batch status command
var batchStatusCmd = &cobra.Command{
	Use:   "status <batch-id>",
	Short: "Get batch job status",
	Long: `Get the status and progress of a batch job.

The batch job ID is typically returned when you initiate a batch operation
(e.g., bulk file upload).

Examples:
  cloud-storage-api-cli batch status 550e8400-e29b-41d4-a716-446655440000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		batchID := args[0]

		// Validate UUID format
		if err := util.ValidateUUID(batchID); err != nil {
			return fmt.Errorf("invalid batch ID: %w", err)
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Get batch job status
		path := fmt.Sprintf("/api/batches/%s/status", batchID)
		var batchResp file.BatchJobResponse
		if err := apiClient.Get(path, &batchResp); err != nil {
			return fmt.Errorf("failed to get batch job status: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(batchResp)
		}

		// Display batch job status
		displayBatchStatus(&batchResp)

		return nil
	},
}

// displayBatchStatus displays batch job status in a formatted way
func displayBatchStatus(batch *file.BatchJobResponse) {
	fmt.Println("\nBatch Job Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Batch ID:        %s\n", batch.BatchID)
	fmt.Printf("Job Type:        %s\n", batch.JobType)
	fmt.Printf("Status:          %s\n", formatBatchStatus(batch.Status))
	fmt.Printf("Progress:        %d%%\n", batch.Progress)
	fmt.Printf("Total Items:     %d\n", batch.TotalItems)
	fmt.Printf("Processed Items: %d\n", batch.ProcessedItems)
	fmt.Printf("Failed Items:    %d\n", batch.FailedItems)

	// Show progress bar
	displayProgressBar(batch.Progress)

	if batch.StartedAt != nil {
		fmt.Printf("Started At:      %s\n", batch.StartedAt.Format(time.RFC3339))
	}
	if batch.EstimatedCompletion != nil {
		fmt.Printf("Estimated Completion: %s\n", batch.EstimatedCompletion.Format(time.RFC3339))
		remaining := time.Until(*batch.EstimatedCompletion)
		if remaining > 0 {
			fmt.Printf("Time Remaining:  %s\n", formatDuration(remaining))
		}
	}
	if batch.ErrorMessage != "" {
		fmt.Printf("\nError: %s\n", batch.ErrorMessage)
	}
	fmt.Println()
}

// formatBatchStatus formats batch status with appropriate styling
func formatBatchStatus(status string) string {
	switch strings.ToUpper(status) {
	case "COMPLETED":
		return "✓ " + status
	case "PROCESSING":
		return "⟳ " + status
	case "FAILED":
		return "✗ " + status
	case "QUEUED":
		return "⏳ " + status
	case "CANCELLED":
		return "⊘ " + status
	default:
		return status
	}
}

// displayProgressBar displays a simple text-based progress bar
func displayProgressBar(progress int) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	barWidth := 30
	filled := (progress * barWidth) / 100
	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	fmt.Printf("Progress Bar:    [%s] %d%%\n", bar, progress)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds > 0 {
			return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
		}
		return fmt.Sprintf("%d minutes", minutes)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}
	return fmt.Sprintf("%d hours", hours)
}

func init() {
	rootCmd.AddCommand(batchCmd)
	batchCmd.AddCommand(batchStatusCmd)
}
