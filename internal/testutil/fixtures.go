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
package testutil

import (
	"time"

	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
)

// SampleFileResponse returns a sample file response for testing
func SampleFileResponse() file.FileResponse {
	folderPath := "/documents"
	return file.FileResponse{
		ID:                  "test-id-123",
		Filename:            "test.txt",
		ContentType:         "text/plain",
		FileSize:            1024,
		FolderPath:          &folderPath,
		CloudinaryUrl:       "http://cloudinary.com/test.txt",
		CloudinarySecureUrl: "https://cloudinary.com/test.txt",
		CreatedAt:           time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:           time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

// SampleFileListResponse returns a sample paginated file list response
func SampleFileListResponse() file.PageResponse {
	return file.PageResponse{
		Content: []file.FileResponse{
			{
				ID:        "file-1",
				Filename:  "test1.txt",
				FileSize:  100,
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:        "file-2",
				Filename:  "test2.txt",
				FileSize:  200,
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		TotalElements:    2,
		TotalPages:       1,
		First:            true,
		Last:             true,
		NumberOfElements: 2,
		Pageable: file.PageableResponse{
			PageNumber: 0,
			PageSize:   20,
		},
	}
}

// SampleErrorResponse returns a sample error response
func SampleErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

// SampleErrorResponseWithDetails returns a sample error response with details
func SampleErrorResponseWithDetails(message, details string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"details": details,
	}
}

// SampleAuthResponse returns a sample authentication response
func SampleAuthResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"token":     token,
		"tokenType": "Bearer",
		"expiresIn": 3600,
	}
}

// SampleFileStatisticsResponse returns a sample file statistics response
func SampleFileStatisticsResponse() file.FileStatisticsResponse {
	return file.FileStatisticsResponse{
		TotalFiles:      10,
		TotalSize:       102400,
		AverageFileSize: 10240,
		StorageUsed:     "100 KB",
		ByContentType: map[string]int64{
			"text/plain": 5,
			"image/jpeg": 5,
		},
		ByFolder: map[string]int64{
			"/documents": 5,
			"/photos":    5,
		},
	}
}
