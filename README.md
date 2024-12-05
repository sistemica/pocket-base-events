# PocketBase with Redis Plugin

A Pocketbase instance with a plugin that provides real-time event synchronization through Redis pub/sub.

## Features

- Real-time event propagation for Create/Update/Delete operations
- Bi-directional sync between PocketBase instances (for example)
- Support for external event processing
- Event source tracking to prevent loops
- Built for multi-architecture (amd64/arm64)

## Installation

```bash
go get github.com/sistemica/pocket-engine
```

## Quick Start

1. Run with Docker:
```bash
docker run -v $(pwd)/pb_data:/app/pb_data -p 8090:8090 \
  -e REDIS_URL=host.docker.internal:6379 \
  ghcr.io/yourusername/pocket-engine
```

2. Or build from source:
```bash
go build
./pocket-engine serve --http=0.0.0.0:8090
```

## Configuration

Environment variables:
- `REDIS_URL`: Redis server address (default: "localhost:6379")
- `--http`: HTTP server address (default: "127.0.0.1:8090")

## Event Structure

```json
{
  "event": "create|update|delete",
  "collection": "collection_name",
  "record": {
    "id": "record_id",
    "field": "value"
  },
  "source": "source_identifier"
}
```

## Redis Channels

- Publisher: `pocketbase:events:publisher`
- Subscriber: `pocketbase:events:receiver`

## Testing Events

Listen to events:
```bash
redis-cli subscribe pocketbase:events:publisher
```

Publish test event:
```bash
redis-cli publish pocketbase:events:receiver '{"event":"create","collection":"test","record":{"field":"hello"},"source":"external"}'
```

## Development

1. Clone repository:
```bash
git clone https://github.com/yourusername/pocket-engine
```

2. Build Docker image:
```bash
docker build -t pocket-engine .
```

3. Run tests:
```bash
go test ./...
```

## Deployment

The project includes GitHub Actions workflow for:
- Automated builds
- Multi-architecture support (amd64/arm64)
- GitHub Container Registry publishing
- Daily builds with latest PocketBase version

## License

MIT License

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Create pull request