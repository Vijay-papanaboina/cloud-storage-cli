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
	ID                  string     `json:"id"`
	Filename            string     `json:"filename"`
	ContentType         string     `json:"contentType"`
	FileSize            int64      `json:"fileSize"`
	FolderPath          *string    `json:"folderPath,omitempty"`
	CloudinaryUrl       string     `json:"cloudinaryUrl"`
	CloudinarySecureUrl string     `json:"cloudinarySecureUrl"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// PageResponse represents a paginated response from the API
type PageResponse struct {
	Content          []FileResponse  `json:"content"`
	Pageable         PageableResponse `json:"pageable"`
	TotalElements    int64           `json:"totalElements"`
	TotalPages       int             `json:"totalPages"`
	First            bool            `json:"first"`
	Last             bool            `json:"last"`
	NumberOfElements int             `json:"numberOfElements"`
}

// PageableResponse represents pagination information
type PageableResponse struct {
	PageNumber int          `json:"pageNumber"`
	PageSize   int          `json:"pageSize"`
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

