
# Orderbook Subscription and Ethereum Smart Contract Fetcher

This repository is a technical test project for building an orderbook subscription system using WebSocket and fetching Ethereum smart contract data. The project is designed to handle cryptocurrency trading pairs and store orderbook data in Redis for quick access.

---

## Features

- Real-time orderbook subscription using WebSocket for Binance and OKX.
- Fetch Ethereum smart contract data, including exchange rates.
- Store orderbook data in Redis for fast retrieval.
- Docker-based deployment for easy setup.

---

## Prerequisites

Before starting, ensure you have the following installed:

- [Go](https://go.dev/dl/) (latest version recommended)
- [Docker](https://www.docker.com/get-started) and Docker Compose
- Git for version control

---

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone <repository-url>
cd <project-directory>
```

Install the required dependencies:

```bash
go mod tidy
```

For running in Docker, see the [Deployment](#deployment) section.

---

## Environment Variables

Create a `.env` file in the root directory by copying `.env.example`:

```bash
cp .env.example .env
```

### Development Environment

#### Database
```plaintext
DB_HOST=localhost
DB_PORT=5432
DB_USER={your_user}
DB_PASSWORD={your_password}
DB_NAME=crypto
```

#### Redis
```plaintext
REDIS_URL=localhost:6379
REDIS_PASSWORD=testdev
```

### Docker Deployment Environment

#### Database
```plaintext
DB_HOST=db
DB_PORT=5432
DB_USER={your_user}
DB_PASSWORD={your_password}
DB_NAME=crypto
```

#### Redis
```plaintext
REDIS_URL=redis:6379
REDIS_PASSWORD=testdev
```

Ensure all environment variables are correctly set before running the project.

---

## Usage

### Running Locally

To run the project in development mode:

```bash
go run cmd/main.go
```

### Running in Docker

To run using Docker:

```bash
docker-compose up
```

The application will be accessible at `http://localhost:<APP_PORT>`. 

If available, API documentation can be found at `/api/v1/docs`.

---

## API Reference

### Get Exchange Rate
Retrieve the Ethereum exchange rate.

```http
GET /api/v1/ethereum/fetch-exchange-rate
```

**Sample Response:**
```json
{
    "data": {
        "infura": {
            "exchange_rate": "1.189324680814567950",
            "last_updated": 1736683266,
            "pair": "Crypto:ALL:APEETH/ETH",
            "pair0": "APEETH",
            "pair1": "ETH",
            "source": "ethereum"
        }
    },
    "message": "Success get exchange rate"
}
```

---

### Orderbook Subscription APIs

#### Add Pair
Add a cryptocurrency trading pair.

```http
POST /api/v1/crypto
```

| Body Field | Type     | Description                         |
|------------|----------|-------------------------------------|
| `pair`     | `string` | **Required**. Example: BTC/USDT     |

**Sample Request:**
```json
{
    "pair": "ETH/USDT"
}
```

**Sample Response:**
```json
{
    "message": "Successfully added crypto data"
}
```

---

#### Get All Pairs
Retrieve all trading pairs.

```http
GET /api/v1/crypto
```

**Sample Response:**
```json
{
    "data": [
        {
            "id": 1,
            "pair": "BTC/USDT"
        }
    ],
    "message": "Successfully retrieved all crypto data",
    "total": 1
}
```

---

#### Delete a Pair
Delete a specific trading pair.

```http
DELETE /api/v1/crypto/:id
```

| Parameter | Type     | Description              |
|-----------|----------|--------------------------|
| `id`      | `int`    | **Required**. Pair ID    |

**Sample Response:**
```json
{
    "message": "Successfully deleted crypto data"
}
```

---

#### Start WebSocket Fetching
Start WebSocket connections for orderbook subscription.

```http
POST /api/v1/crypto/start-websocket
```

**Sample Response:**
```json
{
    "message": "WebSocket connections started"
}
```

---

#### Get Orderbook Data
Retrieve orderbook data from Redis.

```http
GET /api/v1/crypto/orderbook
```

| Query Parameter | Type     | Description                               |
|------------------|----------|-------------------------------------------|
| `key`            | `string` | **Required**. Format: `price:Crypto:pair0/pair1` |

**Sample Response:**
```json
{
    "data": {
        "binance": {
            "ask_price": "94285.30",
            "bid_price": "94285.20",
            "last_updated": 1736687052,
            "mid_price": "94285.25000000",
            "pair": "Crypto:BTC/USDT",
            "pair0": "BTC",
            "pair1": "USDT",
            "source": "binance"
        },
        "okx": {
            "ask_price": "94328.9",
            "bid_price": "94328.8",
            "last_updated": 1736687052,
            "mid_price": "94328.85000000",
            "pair": "Crypto:BTC/USDT",
            "pair0": "BTC",
            "pair1": "USDT",
            "source": "okx"
        }
    },
    "message": "Success get orderbook"
}
```

---

## Deployment

To deploy the application locally using Docker:

```bash
docker-compose up
```

This will start all necessary services, including the API and database.

Stop the containers with:

```bash
docker-compose down
```

The application will be accessible at `http://localhost:<APP_PORT>`.

---

## Contributing

If you’d like to contribute, feel free to fork this repository and submit a pull request. Ensure your code follows the best practices and includes tests.

---

Thank you!
