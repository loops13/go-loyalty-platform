# Architecture & Design Decisions

## Project Structure (Idiomatic Go)

```
cmd/
  server/           # Application entrypoint only
    main.go
internal/
  domain/           # Core business logic & models
    models.go
    errors.go
  dto/              # Request/Response contracts (API surface)
    dto.go
  service/          # Business logic orchestration
    interfaces.go   # ClientService interface
    service.go      # Service implementation
  store/            # Data persistence layer
    interfaces.go   # Repository interface
    inmemory.go     # In-memory implementation
  api/              # HTTP transport layer
    handlers.go     # Endpoint handlers
```

## Architectural Principles

### 1. **Clear Separation of Concerns**
Each package has a single responsibility:
- `domain`: Entities, enums, and domain-level errors
- `dto`: HTTP request/response contracts (marshaling/unmarshaling)
- `service`: Business logic & orchestration
- `store`: Data persistence (currently in-memory, swappable)
- `api`: HTTP transport (handlers, routing, response formatting)

**Rationale**: Makes testing, maintenance, and future database migrations straightforward.

### 2. **Interface-Based Dependencies**
- Handlers depend on `ClientService` interface, not concrete `Service`
- `Service` depends on `Repository` interface, not concrete `InMemoryStore`
- Enables mocking, testing, and swapping implementations

**Rationale**: SOLID Dependency Inversion Principle. Easy to replace in-memory store with PostgreSQL later.

### 3. **Context Propagation**
All business operations accept `context.Context`:
```go
func (s *Service) CreateClient(ctx context.Context, name, email string) (*domain.Client, error)
```

**Rationale**: Cancellation, timeouts, and tracing propagate through the call chain. Standard Go pattern.

### 4. **Domain-First Error Handling**
Custom `DomainError` type with semantic codes:
```go
var (
  ErrClientNotFound = &DomainError{Code: "CLIENT_NOT_FOUND", ...}
  ErrInsufficientPts = &DomainError{Code: "INSUFFICIENT_POINTS", ...}
)
```

**Rationale**: 
- Clients know exactly what went wrong (semantic errors, not strings)
- Handlers can map errors to HTTP status codes consistently
- Testable: `errors.As(err, &de)` for type-safe error checking

### 5. **Slim Service Layer**
Service methods do NOT pass-through; they orchestrate:
- `AwardPoints()` validates award type, looks up points, delegates to repo
- `Redeem()` validates client, fetches reward, checks balance, executes atomically
- `GetAwards()` verifies client exists before returning (business rule enforcement)

**Rationale**: Service is the place for business rules that span multiple entities or need validation.

### 6. **Typed Request/Response DTOs**
Separated from domain models:
- `CreateClientReq`, `AwardPointsReq`, `RedeemReq` for input
- `ClientResp`, `AwardResp`, `RewardResp`, `RedeemResp`, `ErrorResp` for output

**Rationale**:
- Domain models never leak to HTTP layer
- Can evolve API contract independently of domain
- Request/response versions for future versioning (v2 endpoints)
- Cleaner serialization control

### 7. **Handler Factory Pattern**
Each route returns `http.HandlerFunc`, not inline closure:
```go
r.Post("/clients", createClient(svc))

func createClient(svc ClientService) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) { ... }
}
```

**Rationale**: 
- Cleaner route registration (2 lines vs 10+)
- Testable: can call handler directly in unit tests
- Service injected once at registration, not per request
- Easy to add middleware per-route if needed

### 8. **Enum Typing for Award Types**
```go
type AwardType string

const (
  AwardMonthlyContribution AwardType = "MONTHLY_CONTRIBUTION"
  // ...
)

func ValidAwardType(s string) bool
func PointsForAward(t AwardType) int64
```

**Rationale**:
- Type-safe instead of magic strings
- `ValidAwardType()` centralizes validation logic
- `AwardPointsMap` lives in domain, not scattered

### 9. **Validation Layers**
- **HTTP layer**: Non-empty string checks (basic hygiene)
- **Domain layer**: Business rules (award types, point calculations)
- **Service layer**: Orchestration & cross-entity validation (balance checks)

**Rationale**: Each layer is responsible for its concern. Handlers never call domain logic directly.

### 10. **Error Response Formatting**
Standardized error DTO:
```json
{
  "code": "INSUFFICIENT_POINTS",
  "message": "insufficient point balance"
}
```

**Rationale**: Clients parse error codes programmatically, not fragile string matching.

## Removed Anti-Patterns

### ❌ Pass-through Service Methods
**Before:**
```go
func (s *Service) GetClient(id string) (*store.Client, bool) {
  return s.store.GetClient(id)
}
```

**After:** Service only exposes methods with business logic.

**Rationale:** Service layer should add value, not just proxy.

### ❌ Anonymous Request Structs
**Before:**
```go
var req struct{ Type string `json:"type"` }
```

**After:**
```go
var req dto.AwardPointsReq
```

**Rationale**: Named types are discoverable, testable, reusable, documented.

### ❌ Domain Models in Storage Layer
**Before:** `type Client struct` in `store` package

**After:** `type Client struct` in `domain` package

**Rationale:** Models are domain concepts, not storage concepts. Decouples data layer from domain.

### ❌ Store Imported by Handlers
**Before:** Handlers knew about `store.Store`, `store.Client`

**After:** Handlers depend only on `service.ClientService` interface

**Rationale:** Handlers never touch persistence layer. Clean dependency arrow: handlers → service → repo.

### ❌ No Custom Error Types
**Before:**
```go
return Award{}, fmt.Errorf("client not found")
```

**After:**
```go
return nil, domain.ErrClientNotFound
```

**Rationale:** Testable, semantic, HTTP mapable, consistent error handling.

## Future Extensions

This architecture supports:
- **Database migration**: Implement `store.Repository` with SQL backend
- **Transactions**: Add `txn context.Context` parameter to repo methods
- **Caching**: Wrap service methods with caching layer
- **Logging**: Add structured logging middleware in handlers
- **Tracing**: Use context + OpenTelemetry
- **API versioning**: Create `v2` service interface, new DTOs
- **Authorization**: Add middleware to check user permissions
- **Testing**: Mock `Repository` for unit tests, mock `ClientService` for handler tests

## Command to Run

```bash
go run ./cmd/server
```

Services wire together in `main()` with clean dependency injection.
