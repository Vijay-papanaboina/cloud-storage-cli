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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
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
  delete   - Delete a file from cloud storage
  search   - Search files by filename
  info     - Display file storage information`,
}

// fileUploadCmd represents the file upload command
var fileUploadCmd = &cobra.Command{
	Use:   "upload <filepath>",
	Short: "Upload a file to cloud storage",
	Long: `Upload a file to cloud storage with optional folder path and custom filename.

The file will be associated with your authenticated account.
Use Unix-style paths (forward slashes) for folder paths, e.g., /photos/2024.
If --filename is not provided, the original filename will be used.

Examples:
  cloud-storage-api-cli file upload ./document.pdf
  cloud-storage-api-cli file upload ./photo.jpg --folder-path /photos/2024
  cloud-storage-api-cli file upload ./report.pdf --folder-path /documents --filename custom-report.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		folderPath, _ := cmd.Flags().GetString("folder-path")
		filename, _ := cmd.Flags().GetString("filename")

		// Validate folder path if provided
		if folderPath != "" {
			if err := util.ValidatePath(folderPath); err != nil {
				return fmt.Errorf("invalid folder path: %w", err)
			}
		}

		// Validate filename if provided
		if filename != "" {
			if err := util.ValidateFilename(filename); err != nil {
				return fmt.Errorf("invalid filename: %w", err)
			}
		}

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
		if err := apiClient.UploadFile("/api/files/upload", filePath, folderPath, filename, &fileResp); err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(fileResp)
		}

		// Display success message
		fmt.Println("File uploaded successfully!")
		fmt.Printf("File ID: %s\n", fileResp.ID)
		fmt.Printf("Filename: %s\n", fileResp.Filename)
		fmt.Printf("Content Type: %s\n", fileResp.ContentType)
		fmt.Printf("File Size: %s\n", util.FormatFileSize(fileResp.FileSize))
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
		if err := util.ValidatePageNumber(page); err != nil {
			return err
		}
		if err := util.ValidatePageSize(size); err != nil {
			return err
		}
		// Validate folder path if provided
		if folderPath != "" {
			if err := util.ValidatePath(folderPath); err != nil {
				return fmt.Errorf("invalid folder path: %w", err)
			}
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

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(pageResp)
		}

		// Display results
		displayFileList(&pageResp)

		return nil
	},
}

// fileSearchCmd represents the file search command
var fileSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search files by filename",
	Long: `Search files by filename with pagination and optional filtering options.

The search query will match files whose filename contains the query string.

Examples:
  cloud-storage-api-cli file search document
  cloud-storage-api-cli file search photo --page 0 --size 50
  cloud-storage-api-cli file search report --content-type "application/pdf" --folder-path /documents
  cloud-storage-api-cli file search image --page 1 --size 20`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		
		// Validate query is not empty
		if strings.TrimSpace(query) == "" {
			return fmt.Errorf("search query cannot be empty")
		}

		// Get flags
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
		contentType, _ := cmd.Flags().GetString("content-type")
		folderPath, _ := cmd.Flags().GetString("folder-path")

		// Validate pagination parameters
		if err := util.ValidatePageNumber(page); err != nil {
			return err
		}
		if err := util.ValidatePageSize(size); err != nil {
			return err
		}
		// Validate folder path if provided
		if folderPath != "" {
			if err := util.ValidatePath(folderPath); err != nil {
				return fmt.Errorf("invalid folder path: %w", err)
			}
		}

		// Build query parameters
		params := url.Values{}
		params.Set("q", query)
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(size))
		if contentType != "" {
			params.Set("contentType", contentType)
		}
		if folderPath != "" {
			params.Set("folderPath", folderPath)
		}

		// Build URL with query parameters
		path := "/api/files/search?" + params.Encode()

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch search results
		var pageResp file.PageResponse
		if err := apiClient.Get(path, &pageResp); err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(pageResp)
		}

		// Display results
		displayFileList(&pageResp)

		return nil
	},
}

// fileInfoCmd represents the file info command
var fileInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display file storage information",
	Long: `Get information about your file storage including total files, storage used, files by content type, and files by folder.

Examples:
  cloud-storage-api-cli file info`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch file information
		var fileInfo file.FileStatisticsResponse
		if err := apiClient.Get("/api/files/statistics", &fileInfo); err != nil {
			return fmt.Errorf("failed to get file information: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(fileInfo)
		}

		// Display file information
		displayFileInfo(&fileInfo)

		return nil
	},
}

// displayFileInfo displays file information in a formatted way
func displayFileInfo(fileInfo *file.FileStatisticsResponse) {
	fmt.Println("\nFile Storage Information")
	fmt.Println(strings.Repeat("=", 50))

	// Summary section
	fmt.Println("\nSummary:")
	fmt.Printf("  Total Files:      %d\n", fileInfo.TotalFiles)
	fmt.Printf("  Storage Used:     %s\n", fileInfo.StorageUsed)
	fmt.Printf("  Average File Size: %s\n", util.FormatFileSize(fileInfo.AverageFileSize))

	// By content type section
	if len(fileInfo.ByContentType) > 0 {
		fmt.Println("\nBy Content Type:")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("%-30s %s\n", "Content Type", "Count")
		fmt.Println(strings.Repeat("-", 50))

		// Sort content types alphabetically
		contentTypes := make([]string, 0, len(fileInfo.ByContentType))
		for ct := range fileInfo.ByContentType {
			contentTypes = append(contentTypes, ct)
		}
		sort.Strings(contentTypes)

		for _, ct := range contentTypes {
			fmt.Printf("%-30s %d\n", ct, fileInfo.ByContentType[ct])
		}
	} else {
		fmt.Println("\nBy Content Type: None")
	}

	// By folder section
	if len(fileInfo.ByFolder) > 0 {
		fmt.Println("\nBy Folder:")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("%-30s %s\n", "Folder Path", "Count")
		fmt.Println(strings.Repeat("-", 50))

		// Sort folders alphabetically
		folders := make([]string, 0, len(fileInfo.ByFolder))
		for folder := range fileInfo.ByFolder {
			folders = append(folders, folder)
		}
		sort.Strings(folders)

		for _, folder := range folders {
			fmt.Printf("%-30s %d\n", folder, fileInfo.ByFolder[folder])
		}
	} else {
		fmt.Println("\nBy Folder: None")
	}

	fmt.Println()
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
			id, filename, contentType, util.FormatFileSize(f.FileSize), folder, createdAt)
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
	Use:   "download <file-id-or-path>",
	Short: "Download a file from cloud storage",
	Long: `Download a file from cloud storage to your local filesystem.

You can download by:
  - File ID (UUID): 550e8400-e29b-41d4-a716-446655440000
  - Filepath: /photos/2024/image.jpg or document.pdf (for root folder)

The file will be saved to the specified output path, or to the current directory
if no output path is provided. If the output path is a directory, the file will
be saved with its original filename in that directory.

Examples:
  # Download by UUID
  cloud-storage-api-cli file download 550e8400-e29b-41d4-a716-446655440000
  
  # Download by filepath
  cloud-storage-api-cli file download /photos/2024/image.jpg
  cloud-storage-api-cli file download document.pdf
  
  # Download with custom output
  cloud-storage-api-cli file download /documents/report.pdf --output ./downloads/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		outputPath, _ := cmd.Flags().GetString("output")

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Check if identifier is a UUID or filepath
		var finalPath string
		if err := util.ValidateUUID(identifier); err == nil {
			// It's a UUID - use existing download endpoint
			path := fmt.Sprintf("/api/files/%s/download", identifier)
			finalPath, err = apiClient.DownloadFile(path, outputPath)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
		} else {
			// It's a filepath - use new download-by-path endpoint
			// URL encode the filepath
			encodedPath := url.QueryEscape(identifier)
			path := fmt.Sprintf("/api/files/download-by-path?filepath=%s", encodedPath)
			finalPath, err = apiClient.DownloadFile(path, outputPath)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
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
		fmt.Printf("File size: %s\n", util.FormatFileSize(fileInfo.Size()))

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

		// Validate UUID format
		if err := util.ValidateUUID(fileID); err != nil {
			return err
		}
		// Validate filename if provided
		if filename != "" {
			if err := util.ValidateFilename(filename); err != nil {
				return fmt.Errorf("invalid filename: %w", err)
			}
		}
		// Validate folder path if provided
		if folderPath != "" {
			if err := util.ValidatePath(folderPath); err != nil {
				return fmt.Errorf("invalid folder path: %w", err)
			}
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

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(fileResp)
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

		// Validate UUID format
		if err := util.ValidateUUID(fileID); err != nil {
			return err
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

	// Add search subcommand to file command
	fileCmd.AddCommand(fileSearchCmd)

	// Add info subcommand to file command
	fileCmd.AddCommand(fileInfoCmd)

	// Add flags to upload command
	fileUploadCmd.Flags().String("folder-path", "", "Optional folder path (Unix-style, e.g., /photos/2024)")
	fileUploadCmd.Flags().String("filename", "", "Custom filename (optional, defaults to original filename)")

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

	// Add flags to search command
	fileSearchCmd.Flags().Int("page", 0, "Page number (0-indexed, default: 0)")
	fileSearchCmd.Flags().Int("size", 20, "Page size (default: 20, max: 100)")
	fileSearchCmd.Flags().String("content-type", "", "Filter by content type (e.g., image/jpeg)")
	fileSearchCmd.Flags().String("folder-path", "", "Filter by folder path (e.g., /photos/2024)")
}
