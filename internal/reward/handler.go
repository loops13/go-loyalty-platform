package reward

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"awesomeProject/internal/client"
	"awesomeProject/internal/logging"
)

// Handler wraps a service and provides HTTP handlers.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers reward routes with the router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/rewards", h.List)
	r.Post("/clients/{id}/redeem", h.Redeem)
}

// List handles GET /rewards
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rewards, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	resp := make([]RewardResp, len(rewards))
	for i, rw := range rewards {
		resp[i] = rewardToResp(&rw)
	}
	writeJSON(w, http.StatusOK, resp)
}

// Redeem handles POST /clients/{id}/redeem
func (h *Handler) Redeem(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	id := chi.URLParam(r, "id")
	var req RedeemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("invalid redeem request body", "client_id", id, "error", err)
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	reward, balance, err := h.svc.Redeem(r.Context(), id, req.RewardID)
	if err != nil {
		writeRedeemError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, RedeemResp{
		Reward:  rewardToResp(reward),
		Balance: balance,
	})
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

func writeRedeemError(w http.ResponseWriter, err error) {
	// Check for client errors
	var ce *client.ClientError
	if errors.As(err, &ce) {
		status := http.StatusBadRequest
		if ce.Code == client.ErrNotFound.Code {
			status = http.StatusNotFound
		}
		writeError(w, status, ce.Code, ce.Message)
		return
	}

	// Check for reward errors
	var re *RewardError
	if errors.As(err, &re) {
		status := http.StatusBadRequest
		if re.Code == ErrNotFound.Code {
			status = http.StatusNotFound
		}
		writeError(w, status, re.Code, re.Message)
		return
	}

	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func rewardToResp(r *Reward) RewardResp {
	return RewardResp{
		ID:         r.ID,
		Name:       r.Name,
		PointsCost: r.PointsCost,
	}
}
