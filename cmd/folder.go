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
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Folder management commands",
	Long: `Manage folders in cloud storage.

Available commands:
  create - Create a new folder
  list   - List all folders
  delete - Delete an empty folder
  info   - Display folder information (alias: stats)`,
}

// folderCreateCmd represents the folder create command
var folderCreateCmd = &cobra.Command{
	Use:   "create <path>",
	Short: "Create a new folder",
	Long: `Create a new folder in cloud storage.

The folder path must start with '/' and use Unix-style paths.
Folders are virtual - they exist when files are uploaded to that path.

Examples:
  cloud-storage-api-cli folder create /photos/2024
  cloud-storage-api-cli folder create /documents --description "My documents"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		description, _ := cmd.Flags().GetString("description")

		// Validate path starts with /
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("folder path must start with '/'")
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Build create request
		createReq := file.FolderCreateRequest{
			Path: path,
		}
		if description != "" {
			createReq.Description = &description
		}

		// Create folder
		var folderResp file.FolderResponse
		if err := apiClient.Post("/api/folders", createReq, &folderResp); err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}

		// Display success message
		fmt.Println("Folder created successfully!")
		fmt.Printf("Path: %s\n", folderResp.Path)
		if folderResp.Description != nil {
			fmt.Printf("Description: %s\n", *folderResp.Description)
		}
		fmt.Printf("File Count: %d\n", folderResp.FileCount)
		fmt.Printf("Created At: %s\n", folderResp.CreatedAt.Format(time.RFC3339))

		return nil
	},
}

// folderListCmd represents the folder list command
var folderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all folders",
	Long: `List all folders in cloud storage.

You can optionally filter by parent path to list only folders within a specific directory.

Examples:
  cloud-storage-api-cli folder list
  cloud-storage-api-cli folder list --parent-path /photos`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		parentPath, _ := cmd.Flags().GetString("parent-path")

		// Build URL with query parameters
		path := "/api/folders"
		if parentPath != "" {
			params := url.Values{}
			params.Set("parentPath", parentPath)
			path += "?" + params.Encode()
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch folder list
		// Note: API returns List<FolderResponse> which is serialized as JSON array
		var folders []file.FolderResponse
		if err := apiClient.Get(path, &folders); err != nil {
			return fmt.Errorf("failed to list folders: %w", err)
		}

		// Display results
		displayFolderList(folders)

		return nil
	},
}

// folderDeleteCmd represents the folder delete command
var folderDeleteCmd = &cobra.Command{
	Use:   "delete <path>",
	Short: "Delete an empty folder",
	Long: `Delete a folder from cloud storage.

The folder must be empty (no files) to be deleted. This operation cannot be undone.
You will be prompted for confirmation unless the --force flag is used.

Examples:
  cloud-storage-api-cli folder delete /photos/2024
  cloud-storage-api-cli folder delete /photos/2024 --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Validate path starts with /
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("folder path must start with '/'")
		}

		// Prompt for confirmation if not forced
		if !force {
			fmt.Printf("Are you sure you want to delete folder '%s'? This cannot be undone. (y/N): ", path)
			var response string
			fmt.Scanln(&response)
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Delete cancelled.")
				return nil
			}
		}

		// URL encode the path for query parameter
		params := url.Values{}
		params.Set("path", path)
		apiPath := "/api/folders?" + params.Encode()

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Delete folder
		if err := apiClient.Delete(apiPath); err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}

		// Display success message
		fmt.Printf("Folder '%s' deleted successfully.\n", path)

		return nil
	},
}

// folderInfoCmd represents the folder info command
var folderInfoCmd = &cobra.Command{
	Use:     "info <path>",
	Aliases: []string{"stats"},
	Short:   "Display folder information",
	Long: `Get information about a folder including file count, storage used, and breakdowns by content type.

Examples:
  cloud-storage-api-cli folder info /photos/2024
  cloud-storage-api-cli folder stats /photos/2024`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Validate path starts with /
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("folder path must start with '/'")
		}

		// URL encode the path for query parameter
		params := url.Values{}
		params.Set("path", path)
		apiPath := "/api/folders/statistics?" + params.Encode()

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch folder information
		var folderInfo file.FolderStatisticsResponse
		if err := apiClient.Get(apiPath, &folderInfo); err != nil {
			return fmt.Errorf("failed to get folder information: %w", err)
		}

		// Display folder information
		displayFolderInfo(&folderInfo)

		return nil
	},
}

// displayFolderList displays the folder list in a formatted table
func displayFolderList(folders []file.FolderResponse) {
	if len(folders) == 0 {
		fmt.Println("No folders found.")
		return
	}

	// Print header
	fmt.Printf("\nFolders (Total: %d)\n\n", len(folders))

	// Print table header
	fmt.Printf("%-40s %-30s %-10s %-20s\n",
		"Path", "Description", "Files", "Created At")
	fmt.Println(strings.Repeat("-", 100))

	// Print table rows
	for _, f := range folders {
		// Truncate path if too long
		path := f.Path
		if len(path) > 40 {
			path = path[:37] + "..."
		}

		// Format description
		description := "-"
		if f.Description != nil && *f.Description != "" {
			description = *f.Description
			if len(description) > 30 {
				description = description[:27] + "..."
			}
		}

		// Format date
		createdAt := f.CreatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("%-40s %-30s %-10d %-20s\n",
			path, description, f.FileCount, createdAt)
	}

	fmt.Println(strings.Repeat("-", 100))
	fmt.Println()
}

// displayFolderInfo displays folder information in a formatted way
func displayFolderInfo(folderInfo *file.FolderStatisticsResponse) {
	fmt.Println("\nFolder Information")
	fmt.Println(strings.Repeat("=", 50))

	// Summary section
	fmt.Println("\nSummary:")
	fmt.Printf("  Path:             %s\n", folderInfo.Path)
	fmt.Printf("  Total Files:      %d\n", folderInfo.TotalFiles)
	fmt.Printf("  Storage Used:     %s\n", folderInfo.StorageUsed)
	fmt.Printf("  Average File Size: %s\n", formatFileSize(folderInfo.AverageFileSize))
	fmt.Printf("  Created At:       %s\n", folderInfo.CreatedAt.Format(time.RFC3339))

	// By content type section
	if len(folderInfo.ByContentType) > 0 {
		fmt.Println("\nBy Content Type:")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("%-30s %s\n", "Content Type", "Count")
		fmt.Println(strings.Repeat("-", 50))

		// Sort content types alphabetically
		contentTypes := make([]string, 0, len(folderInfo.ByContentType))
		for ct := range folderInfo.ByContentType {
			contentTypes = append(contentTypes, ct)
		}
		sort.Strings(contentTypes)

		for _, ct := range contentTypes {
			fmt.Printf("%-30s %d\n", ct, folderInfo.ByContentType[ct])
		}
	} else {
		fmt.Println("\nBy Content Type: None")
	}

	fmt.Println()
}

func init() {
	// Add folder command to root
	rootCmd.AddCommand(folderCmd)

	// Add create subcommand to folder command
	folderCmd.AddCommand(folderCreateCmd)

	// Add list subcommand to folder command
	folderCmd.AddCommand(folderListCmd)

	// Add delete subcommand to folder command
	folderCmd.AddCommand(folderDeleteCmd)

	// Add info subcommand to folder command (with stats alias)
	folderCmd.AddCommand(folderInfoCmd)

	// Add flags to create command
	folderCreateCmd.Flags().String("description", "", "Optional folder description")

	// Add flags to list command
	folderListCmd.Flags().String("parent-path", "", "Filter by parent path (e.g., /photos)")

	// Add flags to delete command
	folderDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

