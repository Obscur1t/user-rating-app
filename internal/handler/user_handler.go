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

	var userList []model.User
	var err error

	if sortParam == "" {
		userList, err = u.service.GetAll(ctx)
	} else if sortParam == "desc" {
		userList, err = u.service.GetFilteredByRatingDESC(ctx)
	} else if sortParam == "asc" {
		userList, err = u.service.GetFilteredByRatingASC(ctx)
	} else {
		http.Error(w, "query parameter not allowed", http.StatusBadRequest)
		return
	}

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

	if updateUser.Likes != nil {
		if err := u.service.ChangeLikes(ctx, nickname, *updateUser.Likes); err != nil {
			http.Error(w, "invalid change likes", http.StatusInternalServerError)
			return
		}
	}
	if updateUser.Viewers != nil {
		if err := u.service.ChangeViewers(ctx, nickname, *updateUser.Viewers); err != nil {
			http.Error(w, "invalid change viewers", http.StatusInternalServerError)
			return
		}
	}
	if updateUser.Name != nil {
		if err := u.service.ChangeName(ctx, nickname, *updateUser.Name); err != nil {
			http.Error(w, "invalid change name", http.StatusInternalServerError)
			return
		}
	}
	if updateUser.Nickname != nil {
		if err := u.service.ChangeNickname(ctx, nickname, *updateUser.Nickname); err != nil {
			http.Error(w, "invalid change nickname", http.StatusInternalServerError)
			return
		}
	}
	if updateUser.Name == nil && updateUser.Likes == nil && updateUser.Viewers == nil && updateUser.Nickname == nil {
		http.Error(w, "all data is nil", http.StatusBadRequest)
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
