package create

import (
	"context"
	"encoding/json"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"net/http"
)

type UserCreator interface {
	Create(ctx context.Context, user *userEntity.User) (string, error)
}

type Handler struct {
	log         logger.Logger
	userCreator UserCreator
}

func NewHandler(log logger.Logger, userCreator UserCreator) *Handler {
	return &Handler{
		log:         log,
		userCreator: userCreator,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest

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

	input := CreateUserRequestToDomain(&request)

	output, err := h.userCreator.Create(r.Context(), input)
	if err != nil {
		h.log.Error("failed to create user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := NewCreateUserResponse(output)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
