package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"rating/internal/dto/request"
	responsedto "rating/internal/dto/response"
	"rating/internal/model"
	response "rating/internal/transport/http"
	"strconv"
)

type UserService interface {
	CreateUser(ctx context.Context, dto request.UserRequestDTO) error
	GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error)
	GetUser(ctx context.Context, nickname string) (*model.User, error)
	ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error
	Delete(ctx context.Context, nickname string) error
}

type UserHandler struct {
	service UserService
	logger  *slog.Logger
}

func NewUserHandler(service UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  log,
	}
}

func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userRequestDto request.UserRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&userRequestDto); err != nil {
		response.ResponseErr(u.logger, w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.CreateUser(ctx, userRequestDto); err != nil {

		if errors.Is(err, model.ErrInvalidInput) {
			response.ResponseErr(u.logger, w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, model.ErrAlreadyExists) {
			response.ResponseErr(u.logger, w, http.StatusConflict, err.Error())
			return
		}

		response.ResponseErr(u.logger, w, http.StatusInternalServerError, "internal server error")
		return
	}
	statusMsg := map[string]string{"status": "ok"}
	response.ResponseJSON(u.logger, w, http.StatusCreated, statusMsg)
}

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := r.URL.Query()

	sort := param.Get("sort")

	size, err := strconv.Atoi(param.Get("size"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(param.Get("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * size

	params := request.NewPaginationQuery(size, offset, sort)

	userList, totalCount, err := u.service.GetAll(ctx, params)

	if err != nil {
		if errors.Is(err, model.ErrInvalidInput) {
			response.ResponseErr(u.logger, w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrInvalidSort) {
			response.ResponseErr(u.logger, w, http.StatusBadRequest, err.Error())
			return
		}
		response.ResponseErr(u.logger, w, http.StatusInternalServerError, "internal server error")
		return
	}
	data := responsedto.NewPaginatedResponse(userList, totalCount)

	response.ResponseJSON(u.logger, w, http.StatusOK, data)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	user, err := u.service.GetUser(ctx, nickname)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			response.ResponseErr(u.logger, w, http.StatusNotFound, err.Error())
			return
		}
		response.ResponseErr(u.logger, w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.ResponseJSON(u.logger, w, http.StatusOK, user)
}

func (u *UserHandler) ChangeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")
	var updateUser request.UpdateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		response.ResponseErr(u.logger, w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := u.service.ChangeData(ctx, nickname, updateUser); err != nil {
		if errors.Is(err, model.ErrInvalidInput) {
			response.ResponseErr(u.logger, w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrNotFound) {
			response.ResponseErr(u.logger, w, http.StatusNotFound, err.Error())
			return
		}
		response.ResponseErr(u.logger, w, http.StatusInternalServerError, "internal server error")
		return
	}
	responseMsg := map[string]string{"status": "ok"}
	response.ResponseJSON(u.logger, w, http.StatusOK, responseMsg)
}

func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	if err := u.service.Delete(ctx, nickname); err != nil {
		if errors.Is(err, model.ErrInvalidInput) {
			response.ResponseErr(u.logger, w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrNotFound) {
			response.ResponseErr(u.logger, w, http.StatusNotFound, err.Error())
			return
		}
		response.ResponseErr(u.logger, w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
