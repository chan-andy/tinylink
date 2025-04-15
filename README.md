# URL Shortener Service

A modern URL shortener service built with Go, featuring Redis for caching and PostgreSQL for persistent storage.

## Features

- ✨ URL shortening with custom aliases support
- 📊 Real-time click tracking and analytics
- ⚡ Redis caching for fast redirects
- 🔍 URL validation and normalization
- 🔒 Input sanitization and security checks
- 🐳 Docker support for easy deployment

## Tech Stack

- Go 1.21+
- Gin Web Framework
- Redis for caching
- PostgreSQL for persistent storage
- Docker & Docker Compose

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
```

2. Set up the environment:
```bash
cp .env.example .env
# Edit .env with your configurations if needed
```

3. Start the infrastructure using Docker:
```bash
docker compose up -d
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`

## API Reference

### Create Short URL
```http
POST /api/v1/shorten
Content-Type: application/json

Request:
{
    "long_url": "https://example.com/very/long/url",
    "custom_alias": "my-custom-url"  // optional
}

Response:
{
    "short_url": "http://localhost:8080/abc123",
    "original_url": "https://example.com/very/long/url",
    "created_at": "2024-03-20T15:30:00Z"
}
```

### Health Check
```http
GET /health

Response:
{
    "status": "ok",
    "services": {
        "redis": "up",
        "postgres": "up"
    }
}
```

Note: Short URLs (e.g., `http://localhost:8080/abc123`) automatically redirect to their original URLs.

## Development

### Project Structure
```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/              # Private application code
├── .env.example          # Example environment variables
├── .gitignore           # Git ignore rules
├── docker-compose.yml    # Docker compose configuration
├── go.mod               # Go module file
└── README.md            # This file
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| REDIS_URL | Redis connection URL | localhost:6379 |
| POSTGRES_URL | PostgreSQL connection URL | postgres://postgres:postgres@localhost:5432/urlshortener |
| BASE_URL | Base URL for shortened links | http://localhost:8080 |

## Security Considerations

- URLs are validated and sanitized before processing
- Rate limiting implemented using Redis
- All user inputs are properly sanitized
- No sensitive information is exposed in shortened URLs

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [go-redis](https://github.com/redis/go-redis)
- [godotenv](https://github.com/joho/godotenv) 