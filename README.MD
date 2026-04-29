# Mini Lead CRM - Backend Assessment

A RESTful API for a Lead Management CRM built with Go.

## Tech Stack
- **Language:** Go (Golang) - Chosen for high performance, excellent concurrency support, and its great fit for building robust microservices.
- **Framework:** Gin (`gin-gonic/gin`) - A minimalist, high-speed web framework that provides excellent routing and built-in JSON validation, drastically reducing boilerplate.
- **Database:** PostgreSQL (with GORM) - A robust relational database ideal for ensuring data integrity and handling transactional status changes natively.
- **Caching:** Redis - Used to cache read-heavy endpoints (`GET /leads/:id`) to improve response times and reduce database load (Level 3 Bonus).

## How to Run (Docker)

The easiest way to run the application is using Docker Compose. This spins up the Go app, PostgreSQL database, and Redis cache automatically.

1. Clone the repository.
2. Run the following command:
   ```bash
   docker compose up --build
   ```
3. The API will be available at `http://localhost:8080`.

**Seeding Data:** 
Once the app is running, you can seed sample data by running the provided script:
```bash
bash seed.sh
```

## API Documentation
A Postman collection is included in the repository! You can import `mini_lead_crm.postman_collection.json` directly into Postman or Insomnia to quickly test all endpoints.

## Design Decisions

**State Machine Enforcement**
The status transitions (`NEW -> CONTACTED -> QUALIFIED -> CONVERTED` or `LOST`) are strictly enforced in the service layer (`internal/services/lead_service.go`). This isolates business logic from the handlers and the database, making it easily testable.

**Bulk Operations (Level 2) & Partial Success**
Bulk operations are designed to support "partial success". If one lead fails validation (e.g., bad email format), it doesn't break the entire batch. Gin's strict array-wide validation was bypassed in favor of individual manual struct validation within a loop. The API returns a `207 Multi-Status` response with detailed `total`, `successful`, `failed`, and a `results` array.

**Caching Strategy (Level 3)**
Redis caches the `GET /leads/:id` endpoint for 10 minutes. The cache is invalidated automatically when a lead's status is updated or the lead is deleted. If Redis is unavailable, the service degrades gracefully and falls back to querying PostgreSQL directly, ensuring no downtime.

**Trade-offs Made**
- **GORM AutoMigrate:** Used for rapid development, but in a real-world scenario, I would use explicit schema migration tools (like `golang-migrate`) for better tracking and rollbacks.
- **Single Redis Instance:** Used a single Redis container for simplicity instead of a highly available cluster.

**What I'd do differently at scale**
- Implement message queues (e.g. Kafka) for bulk processing, so bulk requests are processed asynchronously in the background.
- Add robust pagination, sorting, and filtering to the `GET /leads` endpoint to prevent fetching massive amounts of data at once.

**Handling Concurrent Status Transitions**
If two users try to update a lead's status at the exact same millisecond, race conditions could occur. To handle this at scale, I would implement **Optimistic Concurrency Control**. I would add a `version` field to the `Lead` model. Any update query would include the version: `UPDATE leads SET status = 'CONTACTED', version = 2 WHERE id = 1 AND version = 1`. If `rows_affected == 0`, it means another process already updated the lead, and the API would return a `409 Conflict` error, prompting the client to refresh and try again.