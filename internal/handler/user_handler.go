package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"rating/internal/dto/request"
	"rating/internal/model"
)

func responseJSON(w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("%v: failed to marshal", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func responseErr(w http.ResponseWriter, status int, message string) {
	responseJSON(w, status, map[string]string{"error": message})
}

type UserService interface {
	CreateUser(ctx context.Context, dto request.UserRequestDTO) error
	GetAll(ctx context.Context, sort string) ([]model.User, error)
	GetUser(ctx context.Context, nickname string) (*model.User, error)
	ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error
	Delete(ctx context.Context, nickname string) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userRequestDto request.UserRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&userRequestDto); err != nil {
		responseErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.CreateUser(ctx, userRequestDto); err != nil {

		if errors.Is(err, model.ErrInvalidInput) {
			responseErr(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, model.ErrAlreadyExists) {
			responseErr(w, http.StatusConflict, err.Error())
			return
		}

		responseErr(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response := map[string]string{"status": "ok"}
	responseJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := r.URL.Query()

	sortParam := param.Get("sort")

	userList, err := u.service.GetAll(ctx, sortParam)

	if err != nil {
		if errors.Is(err, model.ErrInvalidSort) {
			responseErr(w, http.StatusBadRequest, err.Error())
			return
		}
		responseErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responseJSON(w, http.StatusOK, userList)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	user, err := u.service.GetUser(ctx, nickname)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			responseErr(w, http.StatusNotFound, err.Error())
			return
		}
		responseErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responseJSON(w, http.StatusOK, user)
}

func (u *UserHandler) ChangeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")
	var updateUser request.UpdateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		responseErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.ChangeData(ctx, nickname, updateUser); err != nil {
		if errors.Is(err, model.ErrInvalidInput) {
			responseErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrNotFound) {
			responseErr(w, http.StatusNotFound, err.Error())
			return
		}
		responseErr(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response := map[string]string{"status": "ok"}
	responseJSON(w, http.StatusOK, response)
}

func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	if err := u.service.Delete(ctx, nickname); err != nil {
		if errors.Is(err, model.ErrInvalidInput) {
			responseErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrNotFound) {
			responseErr(w, http.StatusNotFound, err.Error())
			return
		}
		responseErr(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
