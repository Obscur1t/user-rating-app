package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"rating/internal/dto/request"
	"rating/internal/model"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
var serverErr = errors.New("internal server error")

type MockUserService struct {
	UserService

	CreateErr error

	GetAllErr error

	GetUserErr error

	ChangeErr error
	DeleteErr error
}

func (m *MockUserService) CreateUser(ctx context.Context, dto request.UserRequestDTO) error {
	return m.CreateErr
}

func (m *MockUserService) GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error) {
	return nil, 0, m.GetAllErr
}

func (m *MockUserService) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	return nil, m.GetUserErr
}

func (m *MockUserService) ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error {
	return m.ChangeErr
}

func (m *MockUserService) Delete(ctx context.Context, nickname string) error {
	return m.DeleteErr
}

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		dto            request.UserRequestDTO
		mockErr        error
		expectedStatus int
	}{
		{
			name: "success",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid input",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:        model.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "already exists",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:        model.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name: "server error",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:        serverErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserService{
				CreateErr: tt.mockErr,
			}

			bodyBytes, _ := json.Marshal(tt.dto)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			handler := NewUserHandler(&mock, discardLogger)
			handler.CreateUserHandler(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		dto            request.PaginationQuery
		mockErr        error
		expectedStatus int
	}{
		{
			name: "success",
			dto: request.PaginationQuery{
				Sort:   "asc",
				Limit:  10,
				Offset: 2,
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid input",
			dto: request.PaginationQuery{
				Sort:   "asc",
				Limit:  10,
				Offset: 2,
			},
			mockErr:        model.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid sort param",
			dto: request.PaginationQuery{
				Sort:   "asc",
				Limit:  10,
				Offset: 2,
			},
			mockErr:        model.ErrInvalidSort,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "server error",
			dto: request.PaginationQuery{
				Sort:   "asc",
				Limit:  10,
				Offset: 2,
			},
			mockErr:        serverErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserService{
				GetAllErr: tt.mockErr,
			}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)

			q := req.URL.Query()

			q.Add("sort", tt.dto.Sort)
			q.Add("size", strconv.Itoa(tt.dto.Limit))
			q.Add("page", strconv.Itoa(tt.dto.Offset))

			req.URL.RawQuery = q.Encode()

			handler := NewUserHandler(&mock, discardLogger)
			handler.GetUsers(rec, req)
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		nickname       string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "success",
			nickname:       "testNick",
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			nickname:       "testNick",
			mockErr:        model.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		}, {
			name:           "server error",
			nickname:       "testNick",
			mockErr:        serverErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserService{
				GetUserErr: tt.mockErr,
			}

			handler := NewUserHandler(&mock, discardLogger)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /users/{nickname}", handler.GetUser)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.nickname, nil)

			mux.ServeHTTP(rec, req)

			require.Equal(t, rec.Code, tt.expectedStatus)
		})
	}
}

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }

func TestUserHandler_ChangeData(t *testing.T) {
	tests := []struct {
		name           string
		nickname       string
		dto            request.UpdateUserDTO
		mockErr        error
		expectedStatus int
	}{
		{
			name:     "success",
			nickname: "testNick",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nick"),
				Viewers:  ptrInt(50),
				Likes:    ptrInt(25),
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid input",
			nickname: "testNick",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nick"),
				Viewers:  ptrInt(50),
				Likes:    ptrInt(25),
			},
			mockErr:        model.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "not found",
			nickname: "testNick",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nick"),
				Viewers:  ptrInt(50),
				Likes:    ptrInt(25),
			},
			mockErr:        model.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "server error",
			nickname: "testNick",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nick"),
				Viewers:  ptrInt(50),
				Likes:    ptrInt(25),
			},
			mockErr:        serverErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserService{
				ChangeErr: tt.mockErr,
			}

			handler := NewUserHandler(&mock, discardLogger)

			mux := http.NewServeMux()
			mux.HandleFunc("PATCH /users/{nickname}", handler.ChangeData)

			bodyBytes, _ := json.Marshal(tt.dto)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/users/"+tt.nickname, bytes.NewReader(bodyBytes))

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		nickname       string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "success",
			nickname:       "testNick",
			mockErr:        nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid input",
			nickname:       "testNick",
			mockErr:        model.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			nickname:       "testNick",
			mockErr:        model.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "server error",
			nickname:       "testNick",
			mockErr:        serverErr,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserService{
				DeleteErr: tt.mockErr,
			}

			handler := NewUserHandler(&mock, discardLogger)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /users/{nickname}", handler.Delete)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.nickname, nil)

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
