# PocketBase Events

A PocketBase distribution with integrated Redis event synchronization capabilities, designed for real-time data synchronization between multiple PocketBase instances or external services.

## Features

- Complete PocketBase server functionality
- Real-time event propagation for Create/Update/Delete operations
- Bi-directional sync between PocketBase instances
- Support for external event processing
- Event source tracking to prevent loops
- Built for multi-architecture (amd64/arm64)

## Quick Start

1. Run with Docker:
```bash
docker run -v $(pwd)/pb_data:/app/pb_data -p 8090:8090 \
  -e REDIS_URL=host.docker.internal:6379 \
  ghcr.io/sistemica/pocketbase-events
```

2. Or build from source:
```bash
git clone https://github.com/sistemica/pocketbase-events
cd pocketbase-events
go build
./pocketbase-events serve --http=0.0.0.0:8090
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

## Use Cases

1. Multi-Instance Synchronization
    - Run multiple PocketBase instances with shared data
    - Automatic propagation of changes across instances
    - Load balancing and high availability setups

2. External Service Integration
    - React to PocketBase changes in external services
    - Trigger PocketBase updates from external systems
    - Build event-driven architectures

3. Data Replication
    - Maintain data copies across different locations
    - Implement backup strategies
    - Create read replicas for better performance

## Development

1. Clone repository:
```bash
git clone https://github.com/sistemica/pocketbase-events
cd pocketbase-events
```

2. Build Docker image:
```bash
docker build -t pocketbase-events .
```

3. Run tests:
- no tests yet -

## Deployment

The project includes GitHub Actions workflow for:
- Automated builds
- Multi-architecture support (amd64/arm64)
- GitHub Container Registry publishing
- Daily builds with latest PocketBase version

### Docker Registry
Images are available at `ghcr.io/sistemica/pocketbase-events`

Tags:
- `latest`: Most recent build
- `sha-xxxxx`: Specific commit builds
- `vX.Y.Z`: Version releases

## Architecture

The system consists of three main components:

1. PocketBase Core
    - Handles all standard PocketBase operations
    - Manages database and authentication
    - Provides REST API and admin interface

2. Redis Event Plugin
    - Listens for PocketBase record changes
    - Publishes events to Redis
    - Subscribes to external events
    - Processes incoming events

3. Redis
    - Acts as message broker
    - Enables pub/sub communication
    - Provides event distribution

## License

MIT License

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Create pull request

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/sistemica/pocketbase-events/issues).