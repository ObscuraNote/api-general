package http

import (
	"net/http"
	"strings"

	"github.com/ObscuraNote/api-general/internal/users/dto"
	"github.com/ObscuraNote/api-general/internal/users/service"
	"github.com/ObscuraNote/api-general/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/philippe-berto/logger"
)

type handler struct {
	log     *logger.Logger
	service service.UserService
}

func Register(router chi.Router, us service.UserService, log logger.Logger) {
	h := &handler{
		log:     &log,
		service: us,
	}

	router.Post("/users", h.CreateUser)
	router.Get("/users/check", h.CheckUserExists)
	router.Put("/users/password", h.UpdatePassword)
	router.Delete("/users", h.DeleteUser)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInput
	if err := utils.ReadBody(r, &input); err != nil {
		_ = utils.Fault(w, http.StatusBadRequest, utils.InvalidBody)

		return
	}

	if err := h.service.CreateUser(input.UserAddress, input.Password); err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "users", "function": "CreateUser"}).
			Error("Failed to create user")

		_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInput
	input.UserAddress, input.Password = getCredentials(r)
	if input.UserAddress == "" || input.Password == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)

		return
	}

	exists, err := h.service.CheckUserExists(input.UserAddress, input.Password)
	if err != nil {
		_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
	}

	if exists {
		w.WriteHeader(http.StatusOK)
	} else {
		_ = utils.Fault(w, http.StatusNotFound, utils.UserNotFound)
	}

}

func (h *handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdatePasswordInput
	if err := utils.ReadBody(r, &input); err != nil {
		_ = utils.Fault(w, http.StatusBadRequest, utils.InvalidBody)
	}

	if input.UserAddress == "" || input.Password == "" || input.NewPassword == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.InvalidBody)
		return
	}

	if err := h.service.UpdatePassword(input.UserAddress, input.Password, input.NewPassword); err != nil {
		if err.Error() == utils.UserNotFound {
			_ = utils.Fault(w, http.StatusUnauthorized, utils.InvalidCredentials)
		} else {
			_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
		}
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserInput
	if err := utils.ReadBody(r, &input); err != nil {
		_ = utils.Fault(w, http.StatusBadRequest, utils.InvalidBody)

		return
	}

	if input.UserAddress == "" || input.Password == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)

		return
	}

	deleted, err := h.service.DeleteUser(input.UserAddress, input.Password)
	if err != nil {
		if err.Error() == utils.UserNotFound {
			_ = utils.Fault(w, http.StatusUnauthorized, utils.InvalidCredentials)
		} else {
			_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
		}
		return
	}

	if !deleted {
		_ = utils.Fault(w, http.StatusNotFound, utils.UserNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getCredentials(r *http.Request) (string, string) {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		credentials := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(credentials, ":")
		userAddress := parts[0]
		password := parts[1]

		return userAddress, password
	}

	return "", ""
}
