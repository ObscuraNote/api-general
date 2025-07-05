package http

import (
	"net/http"
	"strings"

	"github.com/ObscuraNote/api-general/internal/keys/dto"
	kService "github.com/ObscuraNote/api-general/internal/keys/service"
	uService "github.com/ObscuraNote/api-general/internal/users/service"
	"github.com/ObscuraNote/api-general/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/philippe-berto/logger"
)

type handler struct {
	log *logger.Logger
	us  uService.UserService
	ks  kService.KeysService
}

func Register(router chi.Router, ks kService.KeysService, us uService.UserService, log logger.Logger) {
	h := &handler{
		log: &log,
		us:  us,
		ks:  ks,
	}

	router.Post("/keys", h.AddKey)
	router.Get("/keys", h.GetKeysByUser)
	router.Delete("/keys/{id}", h.DeleteKey)
}

func (h *handler) AddKey(w http.ResponseWriter, r *http.Request) {
	var input dto.KeyImput
	if err := utils.ReadBody(r, &input); err != nil {
		_ = utils.Fault(w, http.StatusBadRequest, utils.InvalidBody)
		return
	}

	if input.UserAddress == "" || input.Password == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	createdKey, err := h.ks.AddKey(input)
	if err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "keys", "function": "AddKey"}).
			Error("Failed to add key")

		if err.Error() == utils.ErrUnauthorized {
			_ = utils.Fault(w, http.StatusUnauthorized, utils.InvalidCredentials)
		} else {
			_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
		}
		return
	}

	if err := utils.WriteBody(w, http.StatusCreated, createdKey); err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "keys", "function": "AddKey"}).
			Error("Failed to write response")
		return
	}
}

func (h *handler) GetKeysByUser(w http.ResponseWriter, r *http.Request) {
	userAddress, password := getCredentials(r)
	if userAddress == "" || password == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	auth := dto.AuthInput{
		UserAddress: userAddress,
		Password:    password,
	}

	keys, err := h.ks.GetKeysByUser(r.Context(), auth)
	if err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "keys", "function": "GetKeysByUser"}).
			Error("Failed to get keys")

		if err.Error() == utils.ErrUnauthorized {
			_ = utils.Fault(w, http.StatusUnauthorized, utils.InvalidCredentials)
		} else {
			_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
		}
		return
	}

	if keys == nil {
		keys = []dto.KeyOutput{}
	}

	h.log.WithFields(logger.Fields{"keys": keys, "domain": "keys", "function": "GetKeysByUser"}).Debug("Retrieved keys")

	if err := utils.WriteBody(w, http.StatusOK, keys); err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "keys", "function": "GetKeysByUser"}).
			Error("Failed to write response")
		return
	}
}

func (h *handler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	userAddress, password := getCredentials(r)
	keyID := chi.URLParam(r, "id")
	if keyID == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	auth := dto.AuthInput{
		UserAddress: userAddress,
		Password:    password,
	}

	if auth.UserAddress == "" || auth.Password == "" {
		_ = utils.Fault(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	if err := h.ks.DeleteKey(keyID, auth); err != nil {
		h.log.WithFields(logger.Fields{"error": err.Error(), "domain": "keys", "function": "DeleteKey"}).
			Error("Failed to delete key")

		if err.Error() == utils.ErrUnauthorized {
			_ = utils.Fault(w, http.StatusUnauthorized, utils.InvalidCredentials)
		} else {
			_ = utils.Fault(w, http.StatusInternalServerError, utils.InternalCode)
		}
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
