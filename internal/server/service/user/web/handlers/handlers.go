package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nickzhog/userapi/internal/server/repositories"
	"github.com/nickzhog/userapi/internal/server/service/user"
	"github.com/nickzhog/userapi/pkg/logging"
)

type handler struct {
	logger *logging.Logger
	repositories.Repositories
}

func NewHandler(logger *logging.Logger, reps repositories.Repositories) *handler {
	return &handler{
		logger:       logger,
		Repositories: reps,
	}
}

func (h *handler) GetRouteGroup() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", h.searchUsers)
		r.Post("/", h.createUser)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.getUser)
			r.Patch("/", h.updateUser)
			r.Delete("/", h.deleteUser)
		})
	}
}

func (h *handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repositories.User.FindAll(r.Context())
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var createRequest CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	usr, err := user.NewUser(createRequest.DisplayName, createRequest.Email)
	if err != nil {
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	err = h.Repositories.User.Create(r.Context(), &usr)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

func (h *handler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	usr, err := h.Repositories.User.FindOne(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ErrUserNotFound(err).Render(w, r)
			return
		}
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var usr user.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	err = h.Repositories.User.Update(r.Context(), id, usr)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ErrUserNotFound(err).Render(w, r)
			return
		}
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Repositories.User.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ErrUserNotFound(err).Render(w, r)
			return
		}
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}
}
