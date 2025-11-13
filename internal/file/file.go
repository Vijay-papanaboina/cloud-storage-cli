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
package file

import "time"

// FileResponse represents file information from the API
type FileResponse struct {
	ID                  string    `json:"id"`
	Filename            string    `json:"filename"`
	ContentType         string    `json:"contentType"`
	FileSize            int64     `json:"fileSize"`
	FolderPath          *string   `json:"folderPath,omitempty"`
	CloudinaryUrl       string    `json:"cloudinaryUrl"`
	CloudinarySecureUrl string    `json:"cloudinarySecureUrl"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// PageResponse represents a paginated response from the API
type PageResponse struct {
	Content          []FileResponse   `json:"content"`
	Pageable         PageableResponse `json:"pageable"`
	TotalElements    int64            `json:"totalElements"`
	TotalPages       int              `json:"totalPages"`
	First            bool             `json:"first"`
	Last             bool             `json:"last"`
	NumberOfElements int              `json:"numberOfElements"`
}

// PageableResponse represents pagination information
type PageableResponse struct {
	PageNumber int           `json:"pageNumber"`
	PageSize   int           `json:"pageSize"`
	Sort       *SortResponse `json:"sort,omitempty"`
}

// SortResponse represents sort information
type SortResponse struct {
	Sorted    bool   `json:"sorted"`
	Direction string `json:"direction,omitempty"`
	Property  string `json:"property,omitempty"`
}

// FileUpdateRequest represents a request to update file metadata
type FileUpdateRequest struct {
	Filename   *string `json:"filename,omitempty"`
	FolderPath *string `json:"folderPath,omitempty"`
}

// FileStatisticsResponse represents file statistics from the API
type FileStatisticsResponse struct {
	TotalFiles      int64            `json:"totalFiles"`
	TotalSize       int64            `json:"totalSize"`
	AverageFileSize int64            `json:"averageFileSize"`
	StorageUsed     string           `json:"storageUsed"`
	ByContentType   map[string]int64 `json:"byContentType"`
	ByFolder        map[string]int64 `json:"byFolder"`
}

// FolderCreateRequest represents a request to create a folder
type FolderCreateRequest struct {
	Path        string  `json:"path"`
	Description *string `json:"description,omitempty"`
}

// FolderResponse represents folder information from the API
type FolderResponse struct {
	Path        string    `json:"path"`
	Description *string   `json:"description,omitempty"`
	FileCount   int64     `json:"fileCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

// FolderListResponse represents a list of folders from the API
type FolderListResponse struct {
	Folders []FolderResponse `json:"folders"`
}

// FolderStatisticsResponse represents folder statistics from the API
type FolderStatisticsResponse struct {
	Path            string           `json:"path"`
	TotalFiles      int64            `json:"totalFiles"`
	TotalSize       int64            `json:"totalSize"`
	AverageFileSize int64            `json:"averageFileSize"`
	StorageUsed     string           `json:"storageUsed"`
	ByContentType   map[string]int64 `json:"byContentType"`
	CreatedAt       time.Time        `json:"createdAt"`
}

// FileUrlResponse represents a signed download URL response from the API
type FileUrlResponse struct {
	URL          string    `json:"url"`
	PublicID     string    `json:"publicId"`
	Format       string    `json:"format"`
	ResourceType string    `json:"resourceType"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
