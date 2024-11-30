# Fetch Rewards Receipt Processor

## Prerequisites

- Go 1.17 or later
- Docker (optional, for containerization)
- Git

## Installation

1. **Install Dependencies**

   ```bash
   go mod download
   ```

2. **Generate Swagger Documentation**

   ```bash
   swag init
   ```

## Running the Application locally

To run the application locally, use the following command:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## API

Swagger UI endpoint:

```bash
http://localhost:8080/swagger/index.html
```


### Endpoints

- **POST /receipts/process**: Process a receipt and return a unique ID.

   ```bash
   curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d '{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": 6.49}], "total": 6.49}'
   ```

- **GET /receipts/{id}/points**: Retrieve the points for a given receipt ID.

   ```bash
   curl http://localhost:8080/receipts/1234567890/points
   ```

## Running the application using Docker

1. **Build the Docker Image**

   ```bash
   docker build -t fetch-rewards-app .
   ```

2. **Run the Docker Container**

   ```bash
   docker run -p 8080:8080 fetch-rewards-app
   ```

## Test Cases

```bash
go test ./...
```