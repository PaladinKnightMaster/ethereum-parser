# Ethereum Parser

A simple Ethereum parser API written in Go.

## Features

- Get the current block number
- Subscribe to an address for transaction tracking
- Get transactions for a subscribed address

## Project Structure

The project is organized as follows:

```
ethereum-parser/
|-- cmd/
|   |-- main.go
|-- config/
|   |-- config.go
|-- internal/
|   |-- api/
|   |   |-- api.go
|   |-- parser/
|       |-- parser.go
|       |-- parser_test.go
|-- shared/
|   |-- types.go
|   |-- utils.go
|-- dockerfile
|-- Readme.md
|-- go.mod
|-- go.sum
```

## Configuration

The application can be configured using environment variables:

- `SERVER_PORT`: The port on which the server will run (default: `8080`)
- `RPC_URL`: The Ethereum JSON-RPC URL (default: `https://cloudflare-eth.com`)

## Running Locally

1. Clone the repository:

    ```sh
    git clone https://github.com/your-username/ethereum-parser.git
    cd ethereum-parser
    ```

2. Run the application:

    ```sh
    go run cmd/main.go
    ```

## Docker

You can build and run the application using Docker:

1. Build the Docker image:

    ```sh
    docker build -t ethereum-parser .
    ```

2. Run the Docker container:

    ```sh
    docker run -p 8080:8080 ethereum-parser
    ```

## API Endpoints

- `GET /current_block`: Get the current block number.
- `POST /subscribe`: Subscribe to an address. Body: `{ "address": "0x123" }`
- `GET /transactions?address=0x123`: Get transactions for a subscribed address.

## Testing

Run the tests using:

```sh
go test ./internal/parser/ -v