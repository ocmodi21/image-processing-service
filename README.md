# Image Processing Service

## Description
This service processes images collected from retail stores. It downloads images, calculates the perimeter (2 * (Height + Width)) of each image, and simulates GPU processing with a random sleep time between 0.1 and 0.4 seconds.

The service provides REST APIs for:
- Submitting jobs with image URLs and store IDs
- Checking job status

## Architecture

The service is built with the following components:

- **Job Queue**: Manages multiple concurrent jobs
- **Image Processor**: Downloads and analyzes images
- **Store Manager**: Validates store IDs against master data
- **API Handlers**: Handles HTTP requests and responses

## Assumptions

- The store master data is available in CSV format at runtime
- Image URLs are publicly accessible
- Concurrent job processing is required for performance
- In-memory storage is sufficient for the job processing (no persistence)
- The service should handle network failures gracefully

## Setting Up and Running

### Prerequisites
- Go 1.18 or later
- Docker (optional)

### Using Docker

1. Clone the repository:
```
git clone https://github.com/youruser/imageprocessor.git
cd imageprocessor
```

2. Place the store master data in the root directory as `store_master.csv`

3. Build and run with Docker Compose:
```
docker-compose up
```

The service will be available at http://localhost:8080

### Without Docker

1. Clone the repository:
```
git clone https://github.com/youruser/imageprocessor.git
cd imageprocessor
```

2. Place the store master data in the root directory

3. Build and run:
```
go build -o imageprocessor ./cmd/server
./imageprocessor
```

You can configure the port and store master path with environment variables:
```
PORT=9090 STORE_MASTER_PATH=./data/store_master.csv ./imageprocessor
```

## API Documentation

### Submit Job

Creates a job to process images.

**URL**: `/api/submit/`
**Method**: `POST`
**Request Payload**:
```json
{
   "count": 2,
   "visits": [
      {
         "store_id": "S00339218",
         "image_url": [
            "https://www.gstatic.com/webp/gallery/2.jpg",
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "2023-03-15T14:30:00Z"
      },
      {
         "store_id": "S01408764",
         "image_url": [
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "2023-03-15T15:45:00Z"
      }
   ]
}
```

**Success Response**:
- **Code**: 201 CREATED
- **Content**:
```json
{
   "job_id": "123"
}
```

**Error Response**:
- **Code**: 400 BAD REQUEST
- **Content**:
```json
{
   "error": "count does not match number of visits"
}
```

### Get Job Status

Checks the status of a job.

**URL**: `/api/status?jobid=123`
**Method**: `GET`

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
   "status": "completed",
   "job_id": "123"
}
```

For failed jobs:
```json
{
   "status": "failed",
   "job_id": "123",
   "error": [
      {
         "store_id": "S00339218",
         "error": "store not found"
      }
   ]
}
```

**Error Response**:
- **Code**: 400 BAD REQUEST
- **Content**:
```json
{}
```

## Development Environment

- **Operating System**: Linux (Ubuntu 22.04)
- **IDE**: Visual Studio Code
- **Programming Language**: Go 1.18
- **Libraries**: Standard Go libraries (no external dependencies)
- **Containerization**: Docker, Docker Compose

## Future Improvements

If given more time, I would implement the following improvements:

1. **Persistence**: Implement a database backend to store job results.
2. **Observability**: Add logging, metrics, and tracing.
3. **Scaling**: Implement a more sophisticated worker pool for better concurrency control.
4. **Authentication and Authorization**: Secure the API endpoints.
5. **Rate Limiting**: Prevent abuse of the API.
6. **Caching**: Cache frequently accessed store data.
7. **Health Checks**: Add health check endpoints.
8. **Graceful Shutdown**: Handle shutdown signals properly.
