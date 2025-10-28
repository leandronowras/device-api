# Device API

RESTful API for device management built with Go, DuckDB for persistence, and BDD testing.

### Tech Choices
* Go:
    - Designed for simplicity and efficiency, making it easy to build scalable and performant applications, particularly suitable for microservices architecture such as RESTful APIs. 	
    - Features built-in support for concurrency, enabling the management of multiple devices simultaneously without blocking operations, which enhances responsiveness and performance. 	
    - Strongly typed with garbage collection, minimizing runtime errors, improving code reliability, and simplifying maintenance.

* **Persistence with DuckDB**: 
  - Offers efficient data storage and retrieval, optimized for analytical workloads.
  - Provides a lightweight, serverless database solution that simplifies deployment and maintenance.
  - Supports complex queries and analytics directly within the database, improving performance and reducing data transfer overhead.

* **Behavior-Driven Development (BDD) Testing**: 
  - Promotes collaboration between developers, testers, and non-technical stakeholders by using natural language scenarios.
  - Ensures that the system meets user requirements and behaves as expected through comprehensive testing.
  - Facilitates early detection of issues, reducing development costs and time by addressing problems before they escalate.

### How It Works

The codebase is organized into three layers that keep concerns separated:

- **Domain Layer** (`internal/device`) — Core business logic: device entities, validation rules, and factory methods. No dependencies on HTTP or database.
- **Repository Layer** (`internal/repository`) — Data persistence abstraction with a DuckDB implementation. Swap storage without touching business logic.
- **HTTP Layer** (`internal/http`) — REST API handlers that translate HTTP requests into domain operations and back to JSON responses.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/devices` | Create device |
| GET | `/v1/devices` | List devices (filter: `brand`, `state`; pagination: `page`, `limit`) |
| GET | `/v1/devices/{id}` | Get device by ID |
| PATCH | `/v1/devices/{id}` | Update device |
| DELETE | `/v1/devices/{id}` | Delete device |

### Device States

- `available` (default)
- `in-use`
- `inactive`

## Development

### Prerequisites

- Go 1.23+
- Docker (optional)

### Running Locally

```bash
go run cmd/app/main.go
```

API runs at `http://localhost:8080/v1/devices`

### Running with Docker

```bash
docker-compose up
```

### Testing

```bash
go test ./...
```

BDD tests located in `test/bdd/`

## Configuration

Database: `./devices.db` (DuckDB)

Port: `8080`

## Task status

<!-- TASKMASTER_EXPORT_START -->

| Project Dashboard |  |
| :-                |:-|
| Task Progress     | ███████████████████░ 94% |
| Done | 16 |
| In Progress | 0 |
| Pending | 0 |
| Deferred | 0 |
| Cancelled | 1 |
|-|-|
| Subtask Progress | ░░░░░░░░░░░░░░░░░░░░ 0% |
| Completed | 0 |
| In Progress | 0 |
| Pending | 0 |


| ID | Title | Status | Priority | Dependencies | Complexity |
| :- | :-    | :-     | :-       | :-           | :-         |
| 1 | Feature: Define Device Domain Model | ✓&nbsp;done | medium | None | N/A |
| 2 | Feature: Create a new device (POST /v1/devices) | ✓&nbsp;done | medium | None | N/A |
| 3 | Feature: Fetch a single device (GET /v1/devices/{id}) | ✓&nbsp;done | medium | None | N/A |
| 4 | Feature: List devices with pagination (GET /v1/devices) | ✓&nbsp;done | medium | None | N/A |
| 5 | Feature: Filter devices by brand | ✓&nbsp;done | medium | None | N/A |
| 6 | Feature: Filter devices by state | ✓&nbsp;done | medium | None | N/A |
| 7 | Feature: Fully update a device (PUT /v1/devices/{id}) | ✓&nbsp;done | medium | None | N/A |
| 8 | Feature: Partially update a device (PATCH /v1/devices/{id}) | ✓&nbsp;done | medium | None | N/A |
| 9 | Feature: Delete a device (DELETE /v1/devices/{id}) | ✓&nbsp;done | medium | None | N/A |
| 10 | Feature: Enforce domain validations | ✓&nbsp;done | medium | None | N/A |
| 11 | Feature: Persistence with non in-memory database | ✓&nbsp;done | medium | None | N/A |
| 12 | Feature: API Documentation (OpenAPI + Docs) | ✓&nbsp;done | medium | None | N/A |
| 13 | Feature: Health and Readiness Probes | x&nbsp;cancelled | medium | None | N/A |
| 14 | Feature: Containerization | ✓&nbsp;done | medium | None | N/A |
| 15 | Feature: Reasonable Test Coverage | ✓&nbsp;done | medium | None | N/A |
| 16 | Feature: Repository & Delivery | ✓&nbsp;done | medium | None | N/A |
| 17 | Feature: Go 1.23+ Compliance | ✓&nbsp;done | medium | None | N/A |

<!-- TASKMASTER_EXPORT_END -->


