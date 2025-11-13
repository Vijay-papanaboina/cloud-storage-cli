# Cloud Storage API CLI

A command-line interface (CLI) tool for interacting with the Cloud Storage API. This tool provides a convenient way to manage files and folders from the terminal using API key authentication.

## Features

- **API Key Authentication**: Verify and store API keys for secure authentication
- **File Management**: Upload, download, list, search, update, and delete files
- **Folder Management**: Create, list, delete folders and view folder statistics
- **Configuration Management**: Manage CLI settings and API keys

## Installation

### Prerequisites

- Go 1.24 or later
- Access to the Cloud Storage API
- An API key (generated from the web interface)

### Build from Source

1. Clone the repository:

```bash
git clone https://github.com/vijay-papanaboina/cloud-storage-cli.git
cd cloud-storage-cli
```

2. Build the CLI:

Create a `.env` file in the project root:

```bash
API_URL=http://localhost:8000
```

Then build:

```bash
make build
```

The Makefile automatically reads `API_URL` from your `.env` file and hardcodes it into the binary - just like Vite does!

**Alternative: Direct build commands**

If you prefer not to use Makefile, you can use the long `-ldflags` command:

```bash
go build -ldflags "-X github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config.BuildTimeAPIURL=http://api.example.com" -o cloud-storage-api-cli
```

**Examples with Makefile:**

```bash
# Build with .env file
make build

# Or override .env for one-time builds
API_URL=https://api.production.com make build

# Predefined targets
make build-dev      # Uses http://localhost:8000
make build-staging  # Uses https://api.staging.com
make build-prod     # Uses https://api.production.com
```

3. (Optional) Install globally:

```bash
# On Linux/macOS
sudo mv cloud-storage-api-cli /usr/local/bin/

# On Windows, add to PATH
```

## Configuration

The CLI stores configuration in `~/.cloud-storage-cli/config.yaml`. You can manage configuration using the `config` command or environment variables.

### API URL Configuration

**Important**: The API URL is **hardcoded at compile time** using build flags. It cannot be changed at runtime via environment variables, config files, or command-line flags. Each environment (dev/staging/prod) should have its own compiled binary with the appropriate URL.

### Environment Variables

- `CLOUD_STORAGE_API_KEY`: API key for authentication (can be set at runtime)

**Note**: `CLOUD_STORAGE_API_URL` is **ignored** - the URL is set at compile time only.

### Config File

The config file is automatically created on first use. Sensitive values (API keys) are stored securely with file permissions 0600 (owner read/write only).

The config file stores:

- `api_key`: Your API key (set via `auth login` command)
- `api_url`: Read-only, shows the compile-time URL (cannot be changed)

## Usage

### Authentication

The CLI uses API key authentication. API keys must be generated from the web interface at the Settings page.

#### Login (Verify and Store API Key)

```bash
cloud-storage-api-cli auth login
# API key will be prompted securely
```

This command verifies your API key and saves it to the configuration file for future use.

#### View Current User

```bash
cloud-storage-api-cli auth status
```

Displays information about the currently authenticated user based on the stored API key.

### File Management

#### Upload File

```bash
cloud-storage-api-cli file upload ./document.pdf
cloud-storage-api-cli file upload ./photo.jpg --folder-path /photos/2024
```

#### List Files

```bash
cloud-storage-api-cli file list
cloud-storage-api-cli file list --page 0 --size 50
cloud-storage-api-cli file list --sort "filename,asc" --content-type "image/jpeg"
```

#### Search Files

```bash
cloud-storage-api-cli file search document
cloud-storage-api-cli file search photo --page 0 --size 50
```

#### Download File

```bash
cloud-storage-api-cli file download <file-id>
cloud-storage-api-cli file download <file-id> --output ./downloads/
```

#### Update File

```bash
cloud-storage-api-cli file update <file-id> --filename newname.pdf
cloud-storage-api-cli file update <file-id> --folder-path /documents
```

#### Delete File

```bash
cloud-storage-api-cli file delete <file-id>
cloud-storage-api-cli file delete <file-id> --confirm
```

#### File Statistics

```bash
cloud-storage-api-cli file info
```

### Folder Management

#### Create Folder

```bash
cloud-storage-api-cli folder create /photos/2024
cloud-storage-api-cli folder create /documents --description "My documents"
```

#### List Folders

```bash
cloud-storage-api-cli folder list
cloud-storage-api-cli folder list --parent-path /photos
```

#### Delete Folder

```bash
cloud-storage-api-cli folder delete /photos/2024
cloud-storage-api-cli folder delete /photos/2024 --force
```

#### Folder Information

```bash
cloud-storage-api-cli folder info /photos/2024
cloud-storage-api-cli folder stats /photos/2024  # alias
```

### Configuration

#### Show Configuration

```bash
cloud-storage-api-cli config show
```

#### Get Configuration Value

```bash
cloud-storage-api-cli config get api-url  # Shows compile-time URL (read-only)
cloud-storage-api-cli config get api-key
```

**Note**: The `api-url` value is read-only and shows the URL that was hardcoded at compile time. It cannot be changed via the config command.

## Command-Line Options

### Global Flags

- `--config <path>`: Specify config file path
- `--verbose, -v`: Enable verbose output
- `--json`: Output in JSON format

**Note**: The `--api-url` flag has been removed. The API URL is hardcoded at compile time and cannot be changed at runtime.

### Examples

```bash
# Enable verbose output
cloud-storage-api-cli -v file upload document.pdf

# Output in JSON format
cloud-storage-api-cli --json file list
```

## Input Validation

The CLI validates all inputs to ensure security and correctness:

- **UUIDs**: Validated for proper format (8-4-4-4-12 hex digits)
- **Paths**: Must start with `/`, no path traversal (`..`), no backslashes
- **Filenames**: No path separators, no control characters, no reserved names
- **Pagination**: Page number >= 0, page size 1-100

## Error Handling

All errors include context for debugging:

- HTTP errors include status code, method, and URL
- Network errors include request details
- Validation errors include specific field and reason

Example error output:

```
API error (404) [GET http://localhost:8000/api/files/123]: File not found
```

## Security Features

- **Secure API Key Input**: API keys are never passed as command-line arguments
- **Config File Permissions**: Configuration files use 0600 permissions (owner-only access)
- **Credential Masking**: Sensitive values are masked when displayed
- **Input Sanitization**: All user inputs are validated and sanitized
- **Filename Sanitization**: Downloaded files are sanitized to prevent path traversal

## Examples

### Complete Workflow

```bash
# 1. Login with API key (verify and store)
cloud-storage-api-cli auth login

# 2. Check authentication status
cloud-storage-api-cli auth status

# 3. Create a folder
cloud-storage-api-cli folder create /documents

# 4. Upload files
cloud-storage-api-cli file upload ./report.pdf --folder-path /documents
cloud-storage-api-cli file upload ./photo.jpg --folder-path /photos

# 5. List files
cloud-storage-api-cli file list --folder-path /documents

# 6. Search files
cloud-storage-api-cli file search report

# 7. Download a file
cloud-storage-api-cli file download <file-id> --output ./downloads/

# 8. View statistics
cloud-storage-api-cli file info
cloud-storage-api-cli folder info /documents
```

## Troubleshooting

### Authentication Issues

If you're having authentication issues:

1. Check your API key:

```bash
cloud-storage-api-cli config show
```

2. Verify your API key:

```bash
cloud-storage-api-cli auth status
```

3. Re-login with a new API key:

```bash
cloud-storage-api-cli auth login
```

**Note**: API keys must be generated from the web interface. If you need a new API key, visit the Settings page in the web application.

### Network Issues

If you're experiencing network errors:

1. Verify API URL (shows compile-time URL):

```bash
cloud-storage-api-cli config get api-url
```

2. Test connectivity:

```bash
curl $(cloud-storage-api-cli config get api-url)/health
```

3. Check verbose output:

```bash
cloud-storage-api-cli -v file list
```

**Note**: If the API URL is incorrect, you need to rebuild the CLI with the correct URL using build flags (see Build from Source section).

### Configuration Issues

If configuration isn't working:

1. Check config file location:

```bash
cloud-storage-api-cli config show
```

2. Verify file permissions (should be 0600):

```bash
ls -l ~/.cloud-storage-cli/config.yaml
```

3. Use environment variable for API key:

```bash
export CLOUD_STORAGE_API_KEY=your-api-key-here
cloud-storage-api-cli file list
```

**Note**: The API URL cannot be changed via environment variables - it must be set at compile time.

## Development

### Project Structure

```
cloud-storage-cli/
├── cmd/              # CLI commands
│   ├── auth.go       # Authentication commands (API key verification)
│   ├── file.go       # File management commands
│   ├── folder.go     # Folder management commands
│   ├── config.go     # Configuration commands
│   └── root.go       # Root command
├── internal/
│   ├── client/       # HTTP client
│   ├── config/       # Configuration management
│   ├── file/         # File-related types
│   └── util/         # Utility functions
├── main.go           # Entry point
└── go.mod            # Go module definition
```

### Building

**Recommended: Use `.env` file with Makefile**

1. Create a `.env` file:

```bash
API_URL=http://localhost:8000
```

2. Build:

```bash
make build
```

**Alternative: Direct Go build commands**

```bash
# Build for current platform with default URL (localhost:8000)
go build -o cloud-storage-api-cli

# Build with custom API URL
go build -ldflags "-X github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config.BuildTimeAPIURL=https://api.example.com" -o cloud-storage-api-cli

# Build for specific platforms (read API_URL from .env or use default)
GOOS=linux GOARCH=amd64 make build
GOOS=windows GOARCH=amd64 make build
GOOS=darwin GOARCH=amd64 make build
```

### CI/CD Integration

#### GitHub Actions Example

```yaml
- name: Build CLI
  env:
    API_URL: ${{ secrets.API_URL }}
  run: make build
```

Or without Makefile:

```yaml
- name: Build CLI
  run: |
    go build -ldflags "-X github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config.BuildTimeAPIURL=${{ secrets.API_URL }}" -o cloud-storage-api-cli .
```

#### AWS CodeBuild Example

```yaml
build:
  commands:
    - make build
```

Or without Makefile:

```yaml
build:
  commands:
    - go build -ldflags "-X github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config.BuildTimeAPIURL=${API_URL}" -o cloud-storage-api-cli .
```

## License

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

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go best practices
2. All inputs are validated
3. Error messages are clear and helpful
4. Security best practices are followed
5. Tests are included for new features

## Support

For issues, questions, or contributions, please open an issue on the GitHub repository.
