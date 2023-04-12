package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

type UserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

func (ur *UserResponse) Parse(usr user.User) {
	ur.ID = usr.ID
	ur.CreatedAt = usr.CreatedAt
	ur.DisplayName = usr.DisplayName
	ur.Email = usr.Email
}

// @summary Поиск пользователей
// @description Поиск всех пользователей в системе
// @tags users
// @produce json
// @success 200 {array} UserResponse
// @failure 400 {string} string "Неверный запрос"
// @failure 500 {string} string "Внутренняя ошибка сервера"
// @router /v1/users [get]
func (h *handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repositories.User.FindAll(r.Context())
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}

	var response []UserResponse
	for _, usr := range users {
		var elem UserResponse
		elem.Parse(usr)
		response = append(response, elem)
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

// @summary Создание пользователя
// @description Создание нового пользователя в системе
// @tags users
// @accept json
// @produce json
// @param request body UserRequest true "Данные пользователя"
// @success 200 {object} UserResponse
// @failure 400 {string} string "Неверный запрос"
// @failure 500 {string} string "Внутренняя ошибка сервера"
// @router /v1/users [post]
func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var createRequest UserRequest
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

	var response UserResponse
	response.Parse(usr)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

// @summary Получение пользователя
// @description Получение информации о пользователе по его идентификатору
// @tags users
// @produce json
// @param id path string true "Идентификатор пользователя"
// @success 200 {object} UserResponse
// @failure 400 {string} string "Неверный запрос"
// @failure 404 {string} string "Пользователь не найден"
// @failure 500 {string} string "Внутренняя ошибка сервера"
// @router /v1/users/{id} [get]
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

	var response UserResponse
	response.Parse(usr)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

// @summary Обновление пользователя
// @description Обновление информации о пользователе по его идентификатору
// @tags users
// @accept json
// @produce json
// @param id path string true "Идентификатор пользователя"
// @param request body UserRequest true "Данные пользователя"
// @success 200 {object} UserResponse
// @failure 400 {string} string "Неверный запрос"
// @failure 404 {string} string "Пользователь не найден"
// @failure 500 {string} string "Внутренняя ошибка сервера"
// @router /v1/users/{id} [patch]
func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var usr user.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	err = h.Repositories.User.Update(r.Context(), id, &usr)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ErrUserNotFound(err).Render(w, r)
			return
		}
		h.logger.Error(err)
		ErrInvalidRequest(err).Render(w, r)
	}

	var response UserResponse
	response.Parse(usr)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error(err)
		ErrInternalError(err).Render(w, r)
	}
}

// @summary Удаление пользователя
// @description Удаление пользователя по его идентификатору
// @tags users
// @param id path string true "Идентификатор пользователя"
// @success 204 {string} string "Пользователь удален"
// @failure 400 {string} string "Неверный запрос"
// @failure 404 {string} string "Пользователь не найден"
// @failure 500 {string} string "Внутренняя ошибка сервера"
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
