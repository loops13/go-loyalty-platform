package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"awesomeProject/internal/logging"
)

// Handler wraps a service and provides HTTP handlers.
type Handler struct {
	svc *Service
	// rewardSvc will be injected for redeem operations
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers client routes with the router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/clients", h.Create)
	r.Get("/clients/{id}", h.Get)
	r.Post("/clients/{id}/awards", h.Award)
	r.Get("/clients/{id}/awards", h.GetAwards)
}

// Create handles POST /clients
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	var req CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("invalid client create request body", "error", err)
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		logger.Warn("missing client name")
		writeError(w, http.StatusBadRequest, ErrEmptyName.Code, ErrEmptyName.Message)
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		logger.Warn("missing client email")
		writeError(w, http.StatusBadRequest, ErrEmptyEmail.Code, ErrEmptyEmail.Message)
		return
	}

	client, err := h.svc.Create(r.Context(), req.Name, req.Email)
	if err != nil {
		writeClientError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, clientToResp(client))
}

// Get handles GET /clients/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	client, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, clientToResp(client))
}

// Award handles POST /clients/{id}/awards
func (h *Handler) Award(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req AwardReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("invalid award request body", "client_id", id, "error", err)
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	award, err := h.svc.Award(r.Context(), id, strings.TrimSpace(req.Type))
	if err != nil {
		writeClientError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, awardToResp(award))
}

// GetAwards handles GET /clients/{id}/awards
func (h *Handler) GetAwards(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	awards, err := h.svc.GetAwards(r.Context(), id)
	if err != nil {
		writeClientError(w, err)
		return
	}

	resp := make([]AwardResp, len(awards))
	for i, a := range awards {
		resp[i] = awardToResp(&a)
	}
	writeJSON(w, http.StatusOK, resp)
}

// Helpers

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}

func writeClientError(w http.ResponseWriter, err error) {
	var ce *ClientError
	if errors.As(err, &ce) {
		status := http.StatusBadRequest
		if ce.Code == ErrNotFound.Code {
			status = http.StatusNotFound
		}
		writeError(w, status, ce.Code, ce.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func clientToResp(c *Client) ClientResp {
	return ClientResp{
		ID:           c.ID,
		Name:         c.Name,
		Email:        c.Email,
		PointBalance: c.PointBalance,
	}
}

func awardToResp(a *Award) AwardResp {
	return AwardResp{
		ID:            a.ID,
		ClientID:      a.ClientID,
		Type:          string(a.Type),
		PointsAwarded: a.PointsAwarded,
		CreatedAt:     a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
