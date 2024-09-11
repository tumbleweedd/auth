package delete

import (
	"context"
	"encoding/json"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"net/http"
)

type UserRemover interface {
	Delete(ctx context.Context, userUUID string) error
}

type Handler struct {
	log         logger.Logger
	userRemover UserRemover
}

func NewHandler(log logger.Logger, userRemover UserRemover) *Handler {
	return &Handler{
		log:         log,
		userRemover: userRemover,
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var request DeleteUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.log.Error("failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = request.Validate(); err != nil {
		h.log.Error("failed to validate request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := DeleteUserRequestToDomain(&request)

	err = h.userRemover.Delete(r.Context(), input)
	if err != nil {
		h.log.Error("failed to delete user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
