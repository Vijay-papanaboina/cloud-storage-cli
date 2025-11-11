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
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "File management commands",
	Long: `Manage files in cloud storage.

Available commands:
  upload   - Upload a file to cloud storage
  list     - List files (to be implemented)
  download - Download a file (to be implemented)`,
}

// fileUploadCmd represents the file upload command
var fileUploadCmd = &cobra.Command{
	Use:   "upload <filepath>",
	Short: "Upload a file to cloud storage",
	Long: `Upload a file to cloud storage with optional folder path.

The file will be associated with your authenticated account.
Use Unix-style paths (forward slashes) for folder paths, e.g., /photos/2024.

Examples:
  cloud-storage-api-cli file upload ./document.pdf
  cloud-storage-api-cli file upload ./photo.jpg --folder-path /photos/2024
  cloud-storage-api-cli file upload ./report.pdf --folder-path /documents`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		folderPath, _ := cmd.Flags().GetString("folder-path")

		// Validate file exists and is readable
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", filePath)
			}
			return fmt.Errorf("failed to access file: %w", err)
		}

		// Check if it's a directory
		if fileInfo.IsDir() {
			return fmt.Errorf("path is a directory, not a file: %s", filePath)
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Upload file
		var fileResp file.FileResponse
		if err := apiClient.UploadFile("/api/files/upload", filePath, folderPath, &fileResp); err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		// Display success message
		fmt.Println("File uploaded successfully!")
		fmt.Printf("File ID: %s\n", fileResp.ID)
		fmt.Printf("Filename: %s\n", fileResp.Filename)
		fmt.Printf("Content Type: %s\n", fileResp.ContentType)
		fmt.Printf("File Size: %s\n", formatFileSize(fileResp.FileSize))
		if fileResp.FolderPath != nil {
			fmt.Printf("Folder Path: %s\n", *fileResp.FolderPath)
		}
		fmt.Printf("Cloudinary URL: %s\n", fileResp.CloudinaryUrl)
		fmt.Printf("Cloudinary Secure URL: %s\n", fileResp.CloudinarySecureUrl)
		fmt.Printf("Created At: %s\n", fileResp.CreatedAt.Format(time.RFC3339))

		return nil
	},
}

// formatFileSize formats file size in bytes to human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func init() {
	// Add file command to root
	rootCmd.AddCommand(fileCmd)

	// Add upload subcommand to file command
	fileCmd.AddCommand(fileUploadCmd)

	// Add flags to upload command
	fileUploadCmd.Flags().String("folder-path", "", "Optional folder path (Unix-style, e.g., /photos/2024)")
}

