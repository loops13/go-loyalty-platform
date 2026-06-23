# Code Review Summary: Before & After

## Issue #1: Domain Models in Wrong Package
**Before:**
```go
// internal/store/store.go
type Client struct { ... }
type Award struct { ... }
type Reward struct { ... }
```

**After:**
```go
// internal/domain/models.go
type Client struct { ... }
type Award struct { ... }
type Reward struct { ... }
type AwardType string  // Enum instead of magic strings
```

**Why:** Domain concepts belong in domain package, storage layer is just implementation detail.

---

## Issue #2: Unnecessary Pass-Through Methods
**Before:**
```go
// internal/service/service.go
func (s *Service) GetClient(id string) (*store.Client, bool) {
  return s.store.GetClient(id)  // Just passes through
}

func (s *Service) ListRewards() []store.Reward {
  return s.store.ListRewards()  // Just passes through
}
```

**After:**
```go
// internal/service/service.go
func (s *Service) GetClient(ctx context.Context, id string) (*domain.Client, error) {
  c, err := s.repo.GetClient(ctx, id)
  if err != nil {
    return nil, err
  }
  if c == nil {
    return nil, domain.ErrClientNotFound  // Business logic: enforce client exists
  }
  return c, nil
}

func (s *Service) ListRewards(ctx context.Context) ([]domain.Reward, error) {
  return s.repo.ListRewards(ctx)  // Simple delegation is OK when it adds context
}
```

**Why:** Service layer should have business logic. `GetClient()` now validates existence (business rule).

---

## Issue #3: Anonymous Request Structs
**Before:**
```go
// internal/api/handlers.go (inline)
var req struct{
  Name string `json:"name"`
  Email string `json:"email"`
}
```

**After:**
```go
// internal/dto/dto.go
type CreateClientReq struct {
  Name  string `json:"name"`
  Email string `json:"email"`
}

// internal/api/handlers.go
var req dto.CreateClientReq
```

**Why:** Named types are discoverable, testable, document the API contract, enable schema generation.

---

## Issue #4: Mixed Request/Response DTOs
**Before:**
```go
// internal/api/dto.go (only request)
type CreateClientRequest struct { ... }

// handlers.go (response uses domain type)
json.NewEncoder(w).Encode(client)  // client is *domain.Client
```

**After:**
```go
// internal/dto/dto.go
type CreateClientReq struct { ... }
type ClientResp struct { ... }
type AwardResp struct { ... }
type RedeemResp struct { ... }
type ErrorResp struct { ... }

// handlers.go
json.NewEncoder(w).Encode(domainClientToResp(client))
```

**Why:** 
- Explicit contracts for every endpoint
- Can evolve API independently of domain
- Supports versioning (v2 endpoints with different DTOs)
- Cleaner serialization control

---

## Issue #5: No Custom Error Types
**Before:**
```go
// internal/service/service.go
return nil, fmt.Errorf("unknown award type")
return nil, fmt.Errorf("client not found")
```

**After:**
```go
// internal/domain/errors.go
var (
  ErrClientNotFound   = &DomainError{Code: "CLIENT_NOT_FOUND", Message: "..."}
  ErrInvalidAwardType = &DomainError{Code: "INVALID_AWARD_TYPE", Message: "..."}
  ErrInsufficientPts  = &DomainError{Code: "INSUFFICIENT_POINTS", Message: "..."}
)

// handlers.go
if errors.As(err, &de) {
  status := http.StatusBadRequest
  if de.Code == domain.ErrClientNotFound.Code {
    status = http.StatusNotFound
  }
  writeError(w, status, de.Code, de.Message)
}
```

**Why:**
- Testable with `errors.As()`
- Semantic meaning (code != status code)
- Consistent error response format
- Handlers map domain errors → HTTP status/code consistently

---

## Issue #6: No Centralized Award Type Validation
**Before:**
```go
// internal/service/service.go
var awardPoints = map[string]int64{
  "MONTHLY_CONTRIBUTION": 100,
  ...
}

// internal/api/handlers.go
pts, ok := awardPoints[t]
if !ok {
  http.Error(w, "unknown award type", http.StatusBadRequest)
}
```

**After:**
```go
// internal/domain/models.go
type AwardType string

const (
  AwardMonthlyContribution  AwardType = "MONTHLY_CONTRIBUTION"
  ...
)

var AwardPointsMap = map[AwardType]int64{ ... }

func ValidAwardType(s string) bool {
  _, ok := AwardPointsMap[AwardType(s)]
  return ok
}

// handlers.go
if !domain.ValidAwardType(awardType) {
  writeError(w, http.StatusBadRequest, domain.ErrInvalidAwardType.Code, ...)
}
```

**Why:** Single source of truth for award types. Type-safe enum pattern.

---

## Issue #7: Handlers Importing Store Directly
**Before:**
```go
// internal/api/handlers.go
import "awesomeProject/internal/store"

// handlers know about store.Award, store.Reward
func RegisterRoutes(mux *http.ServeMux, store *store.Store) {
  // handlers had tight coupling to store
}
```

**After:**
```go
// internal/api/handlers.go
import "awesomeProject/internal/service"

func RegisterRoutes(r chi.Router, svc service.ClientService) {
  // handlers only know about service interface
}
```

**Why:** Dependency Inversion. Handlers depend on abstractions (service interface), not concrete storage layer.

---

## Issue #8: Handler Methods Too Long & Inline
**Before:**
```go
r.Post("/clients", func(w http.ResponseWriter, req *http.Request) {
  // 15 lines of code inline
  // validation, decoding, service call, response
})
```

**After:**
```go
r.Post("/clients", createClient(svc))

func createClient(svc service.ClientService) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // Clean, extracted, testable
  }
}
```

**Why:** Extracted handlers are testable, composable, and cleaner route registration.

---

## Issue #9: No Error Response Format
**Before:**
```go
http.Error(w, "error message", http.StatusBadRequest)
// Returns plain text
```

**After:**
```go
writeJSON(w, http.StatusBadRequest, dto.ErrorResp{
  Code:    "INVALID_AWARD_TYPE",
  Message: "unknown award type",
})
// Returns JSON: {"code": "INVALID_AWARD_TYPE", "message": "..."}
```

**Why:** Clients parse `code` field programmatically. Consistent error handling.

---

## Issue #10: No Repository Interface
**Before:**
```go
// Service directly imported store.Store
type Service struct {
  store *store.Store
}

// Hard to mock for testing
// Hard to swap implementations
```

**After:**
```go
// Service depends on interface
type Service struct {
  repo store.Repository
}

// Easy to mock with testify/mock or similar
// Easy to implement SQL backend
```

**Why:** Interface-based design = testable + extensible.

---

## Summary of Improvements

| Metric | Before | After |
|--------|--------|-------|
| Packages | 3 | 5 |
| Interfaces | 0 | 2 (ClientService, Repository) |
| Custom Errors | 0 | 8 semantic errors |
| DTOs | 3 request-only | 6 request + response |
| Testability | Low | High (interfaces + DI) |
| Dependency Clarity | Tangled | Clear hierarchy |
| Code Organization | Store mixing concerns | Separated: domain/dto/service/store/api |

## Files Created/Modified

```
✓ internal/domain/models.go         [NEW] Core domain models & enum
✓ internal/domain/errors.go         [NEW] Semantic error types
✓ internal/dto/dto.go               [NEW] Request/Response contracts
✓ internal/service/interfaces.go    [NEW] ClientService interface
✓ internal/service/service.go       [REFACTORED] With context, business logic
✓ internal/store/interfaces.go      [NEW] Repository interface
✓ internal/store/inmemory.go        [NEW] Concrete in-memory implementation
✓ internal/store/store.go           [REMOVED] Old implementation
✓ internal/api/handlers.go          [REFACTORED] Factory pattern, error mapping
✓ internal/api/dto.go               [REMOVED] Moved to internal/dto
✓ cmd/server/main.go                [UPDATED] Clean DI chain
```

## Testing Examples (Enabled by This Refactor)

```go
// Mock Repository for service tests
type mockRepo struct{ ... }
func (m *mockRepo) CreateClient(ctx, name, email) (*domain.Client, error) { ... }

svc := service.New(mockRepo)
client, err := svc.CreateClient(context.Background(), "Alice", "a@x.com")
assert.NoError(t, err)

// Mock Service for handler tests
type mockService struct{ ... }
handler := createClient(mockService)
rr := httptest.NewRecorder()
handler.ServeHTTP(rr, req)
assert.Equal(t, http.StatusCreated, rr.Code)
```

This architecture is now production-ready and fully testable.
