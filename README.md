# Image Processor Service

A scalable image processing service built with Go using microservices architecture. Allows uploading images, processing them asynchronously (e.g., resizing, applying filters), and retrieving processed results via REST API.

## Description

This project is a scalable image processing service built in Go. It consists of two main components:

- **API Server (app)**: Handles HTTP requests for uploading, retrieving, and deleting images.
- **Worker**: Asynchronously processes images using Kafka message queues.

The service uses PostgreSQL for storing image metadata, MinIO for file storage, and Kafka for processing coordination.

## Features

- ✅ Upload images via REST API
- ✅ Asynchronous image processing (resizing, filters, etc.)
- ✅ Retrieve processed images
- ✅ Delete images
- ✅ Scalable architecture using message queues
- ✅ Object storage (S3-compatible)
- ✅ Docker containerization for local development
- ✅ Database migrations
- ✅ Logging and monitoring

## Architecture

The project follows Clean Architecture principles and is divided into layers:

```
cmd/
├── app/          # API server entry point
└── worker/       # Worker entry point

internal/
├── api/          # HTTP handlers and server
├── infra/        # Infrastructure components (Kafka, MinIO, Worker)
├── models/       # Data structures
├── repository/   # Database repository layer
├── service/      # Business logic layer
└── middlewares/  # HTTP middleware

config/           # Configuration files
migrations/       # Database migrations
```

### Image Processing Flow

1. User uploads image via API
2. Image is stored in MinIO, metadata in PostgreSQL
3. Message is sent to Kafka
4. Worker receives message, processes image
5. Processed image is stored back in MinIO
6. Status is updated in database

## Requirements

- Docker and Docker Compose
- Go 1.25.3+ (for local development)
- Make (optional, for convenience)

## Installation and Running

### Quick Start with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/image-processor.git
   cd image-processor
   ```

2. Create `.env` file based on example:
   ```bash
   cp .env.example .env
   ```
   Fill in environment variables (see Configuration section).

3. Start services:
   ```bash
   make up
   # or
   docker compose up
   ```

4. Service will be available at: http://localhost:8080

### Local Development

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Start infrastructure services (DB, Kafka, MinIO):
   ```bash
   docker compose up db kafka minio migrator kafka-init
   ```

3. Start API server:
   ```bash
   go run cmd/app/main.go
   ```

4. In another terminal, start worker:
   ```bash
   go run cmd/worker/main.go
   ```

## Configuration

### Environment Variables (.env)

```env
# Database
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSL_MODE=

# Kafka
KAFKA_BROKERS=
KAFKA_TOPIC=
KAFKA_GROUP_ID=

# MinIO (S3)
MINIO_HOST=
MINIO_PORT=
MINIO_ROOT_USER=
MINIO_ROOT_PASSWORD=
MINIO_SSL=
S3_BUCKET_NAME=
S3_LOCATION=

# Server
SERVER_PORT=
GIN_MODE=
```

### Configuration File (config/local.yaml)

Main settings in YAML format. Supports overriding via environment variables.

## API

### Upload Image

```http
POST /api/v1/images
Content-Type: multipart/form-data

Form data:
- image: image file
- processing: processing type (e.g., "resize:100x100")
```

**Response:**
```json
{
  "id": 1,
  "status": "uploaded"
}
```

### Get Processed Image

```http
GET /api/v1/images/{id}
```

**Response:** Binary image data

### Delete Image

```http
DELETE /api/v1/images/{id}
```

**Response:**
```json
{
  "message": "Image deleted successfully"
}
```

### Check Status

```http
GET /api/v1/images/{id}/status
```

**Response:**
```json
{
  "id": 1,
  "status": "processed"
}
```

## Development

### Project Structure

- `cmd/`: Application entry points
- `internal/api/`: HTTP handlers and routing
- `internal/infra/`: External service integrations
- `internal/models/`: Data models
- `internal/repository/`: Database operations
- `internal/service/`: Business logic
- `migrations/`: Database migrations
- `config/`: Configuration

### Adding New Image Processing

1. Implement processing function in `internal/infra/handlers/images/process_image.go`
2. Update `ImageKafka` model if needed
3. Test the changes

### Linting and Formatting

```bash
make lint
# or
go vet ./...
golangci-lint run ./...
```

## Testing

### Run Tests

```bash
go test ./...
```

### Integration Tests

```bash
# Start infrastructure
docker compose up -d db kafka minio

# Run tests
go test -tags=integration ./...
```