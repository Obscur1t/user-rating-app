package handler

import (
	"encoding/json"
	"net/http"
	"rating/internal/dto/request"
	"rating/internal/model"
	"rating/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userRequestDto request.UserRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&userRequestDto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := u.service.CreateUser(ctx, userRequestDto); err != nil {
		http.Error(w, "invalid create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"status": "ok"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "invalid encode response", http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := r.URL.Query()

	sortParam := param.Get("sort")
	if sortParam != "" && sortParam != "desc" && sortParam != "asc" {
		http.Error(w, "query parameter not allowed", http.StatusBadRequest)
		return
	}

	var userList []model.User
	var err error

	userList, err = u.service.GetAll(ctx, sortParam)

	if err != nil {
		http.Error(w, "invalid get user's list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, "invalid encode response", http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	user, err := u.service.GetUser(ctx, nickname)
	if err != nil {
		http.Error(w, "invalid get user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "invalid encode user", http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) ChangeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")
	var updateUser request.UpdateUserDTO

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		http.Error(w, "invalid decode json", http.StatusBadRequest)
		return
	}

	if err := u.service.ChangeData(ctx, nickname, updateUser); err != nil {
		http.Error(w, "invalid change data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "invalid encode response", http.StatusInternalServerError)
		return
	}

}

func (u *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nickname := r.PathValue("nickname")

	if err := u.service.Delete(ctx, nickname); err != nil {
		http.Error(w, "invalid delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
