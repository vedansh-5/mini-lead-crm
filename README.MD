# Mini Lead CRM 

A RESTful API for a Lead Management CRM built with Go.

## Tech Stack
- **Language:** Go (Golang) - Chosen for high performance, standard concurrency, and its excellent fit for microservices.
- **Framework:** Gin (`gin-gonic/gin`) - A minimalist, high-speed web framework that provides excellent routing and built-in JSON validation.
- **Database:** PostgreSQL (with GORM) - A robust relational database ideal for ensuring data integrity and handling transactional status changes.
- **Caching:** Redis - Used to cache read-heavy endpoints (`GET /leads/:id`) to improve response times (Level 3 Bonus).

## How to Run (Docker)

The easiest way to run the application is using Docker Compose. This spins up the Go app, PostgreSQL database, and Redis cache automatically.

1. Clone the repository.
2. Ensure Docker and Docker Compose are installed.
3. Run the following command:
   ```bash
   docker-compose up --build
4. The API will be available at ```http://localhost:8080```.

## Design Decisions

- State Machine Enforcement: The status transitions (NEW -> CONTACTED -> QUALIFIED -> CONVERTED or LOST) are strictly enforced in the service layer using a helper function. This isolates business logic away from the handlers and the database.

- Bulk Operations (Level 2): Designed to support "partial success". If one lead fails validation, it doesn't break the entire batch. A custom response format returns total, successful, failed, and a detailed results array.

- Caching Strategy (Level 3): Redis is used to cache GET /leads/:id. The cache is invalidated automatically when a lead's status is updated. If Redis fails to connect, the system falls back gracefully to querying PostgreSQL directly, ensuring no downtime.