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
	"net/url"
	"os"
	"strconv"
	"strings"
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
  list     - List files with pagination and filtering
  download - Download a file from cloud storage
  update   - Update file metadata (filename, folder path)
  delete   - Delete a file from cloud storage`,
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

// fileListCmd represents the file list command
var fileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in cloud storage",
	Long: `List files in cloud storage with pagination, sorting, and filtering options.

Examples:
  cloud-storage-api-cli file list
  cloud-storage-api-cli file list --page 0 --size 50
  cloud-storage-api-cli file list --sort "filename,asc"
  cloud-storage-api-cli file list --content-type "image/jpeg" --folder-path /photos
  cloud-storage-api-cli file list --page 1 --size 20 --sort "createdAt,desc"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
		sort, _ := cmd.Flags().GetString("sort")
		contentType, _ := cmd.Flags().GetString("content-type")
		folderPath, _ := cmd.Flags().GetString("folder-path")

		// Validate pagination parameters
		if page < 0 {
			return fmt.Errorf("page must be >= 0")
		}
		if size <= 0 || size > 100 {
			return fmt.Errorf("size must be between 1 and 100")
		}

		// Build query parameters
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(size))
		if sort != "" {
			params.Set("sort", sort)
		}
		if contentType != "" {
			params.Set("contentType", contentType)
		}
		if folderPath != "" {
			params.Set("folderPath", folderPath)
		}

		// Build URL with query parameters
		path := "/api/files"
		if len(params) > 0 {
			path += "?" + params.Encode()
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch file list
		var pageResp file.PageResponse
		if err := apiClient.Get(path, &pageResp); err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		// Display results
		displayFileList(&pageResp)

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

// displayFileList displays the file list in a formatted table
func displayFileList(pageResp *file.PageResponse) {
	if len(pageResp.Content) == 0 {
		fmt.Println("No files found.")
		return
	}

	// Print header
	fmt.Printf("\nFiles (Page %d of %d, Total: %d)\n\n", 
		pageResp.Pageable.PageNumber+1, 
		pageResp.TotalPages, 
		pageResp.TotalElements)

	// Print table header
	fmt.Printf("%-36s %-30s %-20s %-12s %-20s %-20s\n",
		"ID", "Filename", "Content Type", "Size", "Folder", "Created At")
	fmt.Println(strings.Repeat("-", 140))

	// Print table rows
	for _, f := range pageResp.Content {
		// Truncate ID to 36 chars (UUID length)
		id := f.ID
		if len(id) > 36 {
			id = id[:36]
		}

		// Truncate filename if too long
		filename := f.Filename
		if len(filename) > 30 {
			filename = filename[:27] + "..."
		}

		// Truncate content type if too long
		contentType := f.ContentType
		if len(contentType) > 20 {
			contentType = contentType[:17] + "..."
		}

		// Format folder path
		folder := "-"
		if f.FolderPath != nil && *f.FolderPath != "" {
			folder = *f.FolderPath
			if len(folder) > 20 {
				folder = folder[:17] + "..."
			}
		}

		// Format date
		createdAt := f.CreatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("%-36s %-30s %-20s %-12s %-20s %-20s\n",
			id, filename, contentType, formatFileSize(f.FileSize), folder, createdAt)
	}

	// Print pagination info
	fmt.Println(strings.Repeat("-", 140))
	fmt.Printf("Showing %d of %d files", pageResp.NumberOfElements, pageResp.TotalElements)
	if !pageResp.First || !pageResp.Last {
		fmt.Print(" (")
		if !pageResp.First {
			fmt.Print("not first")
		}
		if !pageResp.First && !pageResp.Last {
			fmt.Print(", ")
		}
		if !pageResp.Last {
			fmt.Print("not last")
		}
		fmt.Print(")")
	}
	fmt.Println()
}

// fileDownloadCmd represents the file download command
var fileDownloadCmd = &cobra.Command{
	Use:   "download <file-id>",
	Short: "Download a file from cloud storage",
	Long: `Download a file from cloud storage to your local filesystem.

The file will be saved to the specified output path, or to the current directory
if no output path is provided. If the output path is a directory, the file will
be saved with its original filename in that directory.

Examples:
  cloud-storage-api-cli file download 550e8400-e29b-41d4-a716-446655440000
  cloud-storage-api-cli file download 550e8400-e29b-41d4-a716-446655440000 --output ./downloads/
  cloud-storage-api-cli file download 550e8400-e29b-41d4-a716-446655440000 --output ./myfile.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileID := args[0]
		outputPath, _ := cmd.Flags().GetString("output")

		// Basic UUID format validation (simplified)
		if len(fileID) < 8 {
			return fmt.Errorf("invalid file ID format: %s", fileID)
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Download file
		path := fmt.Sprintf("/api/files/%s/download", fileID)
		finalPath, err := apiClient.DownloadFile(path, outputPath)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		// Get file info for display
		fileInfo, err := os.Stat(finalPath)
		if err != nil {
			// File was downloaded but we can't get info - still show success
			fmt.Printf("File downloaded successfully to: %s\n", finalPath)
			return nil
		}

		// Display success message
		fmt.Println("File downloaded successfully!")
		fmt.Printf("File path: %s\n", finalPath)
		fmt.Printf("File size: %s\n", formatFileSize(fileInfo.Size()))

		return nil
	},
}

// fileUpdateCmd represents the file update command
var fileUpdateCmd = &cobra.Command{
	Use:   "update <file-id>",
	Short: "Update file metadata",
	Long: `Update file metadata (filename and/or folder path).

At least one of --filename or --folder-path must be provided.

Examples:
  cloud-storage-api-cli file update 550e8400-e29b-41d4-a716-446655440000 --filename newname.pdf
  cloud-storage-api-cli file update 550e8400-e29b-41d4-a716-446655440000 --folder-path /documents
  cloud-storage-api-cli file update 550e8400-e29b-41d4-a716-446655440000 --filename newname.pdf --folder-path /documents`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileID := args[0]
		filename, _ := cmd.Flags().GetString("filename")
		folderPath, _ := cmd.Flags().GetString("folder-path")

		// Validate that at least one field is provided
		if filename == "" && folderPath == "" {
			return fmt.Errorf("at least one of --filename or --folder-path must be provided")
		}

		// Basic UUID format validation
		if len(fileID) < 8 {
			return fmt.Errorf("invalid file ID format: %s", fileID)
		}

		// Build update request
		updateReq := file.FileUpdateRequest{}
		if filename != "" {
			updateReq.Filename = &filename
		}
		if folderPath != "" {
			updateReq.FolderPath = &folderPath
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Update file
		path := fmt.Sprintf("/api/files/%s", fileID)
		var fileResp file.FileResponse
		if err := apiClient.Put(path, updateReq, &fileResp); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		// Display success message
		fmt.Println("File updated successfully!")
		fmt.Printf("File ID: %s\n", fileResp.ID)
		fmt.Printf("Filename: %s\n", fileResp.Filename)
		if fileResp.FolderPath != nil {
			fmt.Printf("Folder Path: %s\n", *fileResp.FolderPath)
		} else {
			fmt.Println("Folder Path: (none)")
		}
		fmt.Printf("Updated At: %s\n", fileResp.UpdatedAt.Format(time.RFC3339))

		return nil
	},
}

// fileDeleteCmd represents the file delete command
var fileDeleteCmd = &cobra.Command{
	Use:   "delete <file-id>",
	Short: "Delete a file from cloud storage",
	Long: `Delete a file from cloud storage.

This operation cannot be undone. You will be prompted for confirmation unless
the --confirm flag is used.

Examples:
  cloud-storage-api-cli file delete 550e8400-e29b-41d4-a716-446655440000
  cloud-storage-api-cli file delete 550e8400-e29b-41d4-a716-446655440000 --confirm`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		// Basic UUID format validation
		if len(fileID) < 8 {
			return fmt.Errorf("invalid file ID format: %s", fileID)
		}

		// Prompt for confirmation if not already confirmed
		if !confirm {
			fmt.Printf("Are you sure you want to delete file %s? This cannot be undone. (y/N): ", fileID)
			var response string
			fmt.Scanln(&response)
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Delete file
		path := fmt.Sprintf("/api/files/%s", fileID)
		if err := apiClient.Delete(path); err != nil {
			return fmt.Errorf("delete failed: %w", err)
		}

		// Display success message
		fmt.Printf("File %s deleted successfully.\n", fileID)

		return nil
	},
}

func init() {
	// Add file command to root
	rootCmd.AddCommand(fileCmd)

	// Add upload subcommand to file command
	fileCmd.AddCommand(fileUploadCmd)

	// Add list subcommand to file command
	fileCmd.AddCommand(fileListCmd)

	// Add download subcommand to file command
	fileCmd.AddCommand(fileDownloadCmd)

	// Add update subcommand to file command
	fileCmd.AddCommand(fileUpdateCmd)

	// Add delete subcommand to file command
	fileCmd.AddCommand(fileDeleteCmd)

	// Add flags to upload command
	fileUploadCmd.Flags().String("folder-path", "", "Optional folder path (Unix-style, e.g., /photos/2024)")

	// Add flags to list command
	fileListCmd.Flags().Int("page", 0, "Page number (0-indexed, default: 0)")
	fileListCmd.Flags().Int("size", 20, "Page size (default: 20, max: 100)")
	fileListCmd.Flags().String("sort", "createdAt,desc", "Sort field and direction (e.g., createdAt,desc)")
	fileListCmd.Flags().String("content-type", "", "Filter by content type (e.g., image/jpeg)")
	fileListCmd.Flags().String("folder-path", "", "Filter by folder path (e.g., /photos/2024)")

	// Add flags to download command
	fileDownloadCmd.Flags().StringP("output", "o", "", "Output file path or directory (default: current directory)")

	// Add flags to update command
	fileUpdateCmd.Flags().String("filename", "", "New filename")
	fileUpdateCmd.Flags().String("folder-path", "", "New folder path (Unix-style, e.g., /photos/2024)")

	// Add flags to delete command
	fileDeleteCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
}

