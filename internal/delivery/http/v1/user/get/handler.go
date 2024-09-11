package get

import (
	"context"
	"encoding/json"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"net/http"
)

type UserGetter interface {
	Get(ctx context.Context, uuid string) (*userEntity.User, error)
	GetByUUIDs(ctx context.Context, uuids []string) ([]*userEntity.User, error)
}

type Handler struct {
	log        logger.Logger
	userGetter UserGetter
}

func NewHandler(log logger.Logger, userGetter UserGetter) *Handler {
	return &Handler{
		log:        log,
		userGetter: userGetter,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var request GetUserRequest

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

	input := GetUserRequestToDomain(&request)

	output, err := h.userGetter.Get(r.Context(), input)
	if err != nil {
		h.log.Error("failed to get user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := NewGetUserResponse(output)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.log.Error("failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByUUIDs(w http.ResponseWriter, r *http.Request) {
	var request GetUsersRequest

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

	input := GetUsersRequestToDomain(&request)

	output, err := h.userGetter.GetByUUIDs(r.Context(), input)
	if err != nil {
		h.log.Error("failed to get users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := NewGetUsersResponse(output)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.log.Error("failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
