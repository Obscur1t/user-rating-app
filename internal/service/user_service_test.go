package service

import (
	"context"
	"errors"
	"rating/internal/dto/request"
	"rating/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockUserStore struct {
	UserStore
	CreateErr error

	GetAllErr    error
	GetAllTotal  int
	GetAllResult []model.User

	GetUserErr    error
	GetUserResult *model.User

	ChangeErr error
	DeleteErr error
}

func (m *MockUserStore) Create(ctx context.Context, user model.User) error {
	return m.CreateErr
}

func (m *MockUserStore) GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error) {
	return m.GetAllResult, m.GetAllTotal, m.GetAllErr
}

func (m *MockUserStore) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	return m.GetUserResult, m.GetUserErr
}

func (m *MockUserStore) ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error {
	return m.ChangeErr
}

func (m *MockUserStore) Delete(ctx context.Context, nickname string) error {
	return m.DeleteErr
}

func TestUserService_CreateUser(t *testing.T) {
	serverErr := errors.New("internal server error")
	tests := []struct {
		name        string
		dto         request.UserRequestDTO
		mockErr     error
		expectedErr error
	}{
		{
			name: "success",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name: "empty name",
			dto: request.UserRequestDTO{
				Name:     "",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name: "empty nickname",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name: "negative viewers",
			dto: request.UserRequestDTO{
				Name:     "",
				Nickname: "nickname",
				Viewers:  -100,
				Likes:    50,
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name: "negative likes",
			dto: request.UserRequestDTO{
				Name:     "",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    -50,
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name: "likes > viewers",
			dto: request.UserRequestDTO{
				Name:     "",
				Nickname: "nickname",
				Viewers:  50,
				Likes:    100,
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name: "already exists",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:     model.ErrAlreadyExists,
			expectedErr: model.ErrAlreadyExists,
		},
		{
			name: "server error",
			dto: request.UserRequestDTO{
				Name:     "name",
				Nickname: "nickname",
				Viewers:  100,
				Likes:    50,
			},
			mockErr:     serverErr,
			expectedErr: serverErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserStore{
				CreateErr: tt.mockErr,
			}

			service := NewUserService(&mock)
			err := service.CreateUser(context.Background(), tt.dto)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	serverErr := errors.New("internal server error")
	usersList := []model.User{*model.NewUser("name1", "nickname1", 50, 100), *model.NewUser("name2", "nickname2", 20, 80)}
	tests := []struct {
		name           string
		params         request.PaginationQuery
		mockErr        error
		mockResult     []model.User
		mockTotal      int
		expectedErr    error
		expectedResult []model.User
		expectedTotal  int
	}{
		{
			name:           "success",
			params:         request.NewPaginationQuery(5, 10, "asc"),
			mockErr:        nil,
			mockResult:     usersList,
			mockTotal:      50,
			expectedErr:    nil,
			expectedResult: usersList,
			expectedTotal:  50,
		},
		{
			name:           "unknown sort param",
			params:         request.NewPaginationQuery(5, 10, "sdsds"),
			mockErr:        nil,
			mockResult:     nil,
			mockTotal:      -1,
			expectedErr:    model.ErrInvalidSort,
			expectedResult: nil,
			expectedTotal:  -1,
		},
		{
			name:           "negative limit",
			params:         request.NewPaginationQuery(-5, 10, "asc"),
			mockErr:        nil,
			mockResult:     nil,
			mockTotal:      -1,
			expectedErr:    model.ErrInvalidInput,
			expectedResult: nil,
			expectedTotal:  -1,
		},
		{
			name:           "negative offset",
			params:         request.NewPaginationQuery(5, -10, "asc"),
			mockErr:        nil,
			mockResult:     nil,
			mockTotal:      -1,
			expectedErr:    model.ErrInvalidInput,
			expectedResult: nil,
			expectedTotal:  -1,
		},
		{
			name:           "server error",
			params:         request.NewPaginationQuery(5, 10, "asc"),
			mockErr:        serverErr,
			mockResult:     nil,
			mockTotal:      -1,
			expectedErr:    serverErr,
			expectedResult: nil,
			expectedTotal:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserStore{
				GetAllErr:    tt.mockErr,
				GetAllResult: tt.mockResult,
				GetAllTotal:  tt.mockTotal,
			}

			service := NewUserService(&mock)
			userList, total, err := service.GetAll(context.Background(), tt.params)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, userList, tt.expectedResult)
			require.Equal(t, total, tt.expectedTotal)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	serverErr := errors.New("internal server error")
	user := model.NewUser("name", "nickname", 50, 100)
	test := []struct {
		name           string
		nickname       string
		mockErr        error
		mockResult     *model.User
		expectedErr    error
		expectedResult *model.User
	}{
		{
			name:           "success",
			nickname:       "nickname",
			mockErr:        nil,
			mockResult:     user,
			expectedErr:    nil,
			expectedResult: user,
		},
		{
			name:           "empty nickname",
			nickname:       "",
			mockErr:        nil,
			mockResult:     nil,
			expectedErr:    model.ErrInvalidInput,
			expectedResult: nil,
		},
		{
			name:           "server error",
			nickname:       "nickname",
			mockErr:        serverErr,
			mockResult:     nil,
			expectedErr:    serverErr,
			expectedResult: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserStore{
				GetUserErr:    tt.mockErr,
				GetUserResult: tt.mockResult,
			}

			service := NewUserService(&mock)
			user, err := service.GetUser(context.Background(), tt.nickname)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, user, tt.expectedResult)
		})
	}
}

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }

func TestUserService_Change(t *testing.T) {
	serverErr := errors.New("internal server error")
	tests := []struct {
		name        string
		nickname    string
		dto         request.UpdateUserDTO
		mockErr     error
		expectedErr error
	}{
		{
			name:     "success",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name:     "name is empty",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString(""),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "dto nickname is empty",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString(""),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "likes negative",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(-50),
				Viewers:  ptrInt(100),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "viewers negative",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(-100),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "likes > viewers",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(100),
				Viewers:  ptrInt(50),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "nickname is empty",
			nickname: "",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:        "empty dto",
			nickname:    "nickname",
			dto:         request.UpdateUserDTO{},
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:     "server error",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     serverErr,
			expectedErr: serverErr,
		},
		{
			name:     "server not found",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Name:     ptrString("name"),
				Nickname: ptrString("nickname"),
				Likes:    ptrInt(50),
				Viewers:  ptrInt(100),
			},
			mockErr:     model.ErrNotFound,
			expectedErr: model.ErrNotFound,
		},
		{
			name:     "server invalid input",
			nickname: "nickname",
			dto: request.UpdateUserDTO{
				Likes: ptrInt(160),
			},
			mockErr:     model.ErrInvalidInput,
			expectedErr: model.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserStore{
				ChangeErr: tt.mockErr,
			}

			service := NewUserService(&mock)
			err := service.ChangeData(context.Background(), tt.nickname, tt.dto)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	serverErr := errors.New("server error")
	tests := []struct {
		name        string
		nickname    string
		mockErr     error
		expectedErr error
	}{
		{
			name:        "success",
			nickname:    "nickname",
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "nickname is empty",
			nickname:    "",
			mockErr:     nil,
			expectedErr: model.ErrInvalidInput,
		},
		{
			name:        "server error",
			nickname:    "nickname",
			mockErr:     serverErr,
			expectedErr: serverErr,
		},
		{
			name:        "server not found",
			nickname:    "nickname",
			mockErr:     model.ErrNotFound,
			expectedErr: model.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockUserStore{
				DeleteErr: tt.mockErr,
			}

			service := NewUserService(&mock)
			err := service.Delete(context.Background(), tt.nickname)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
