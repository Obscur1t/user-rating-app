package service

import (
	"context"
	"fmt"
	"rating/internal/dto/request"
	"rating/internal/model"
)

type UserStore interface {
	Create(ctx context.Context, user model.User) error
	GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error)
	GetUser(ctx context.Context, nickname string) (*model.User, error)
	ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error
	Delete(ctx context.Context, nickname string) error
}

type UserService struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, dto request.UserRequestDTO) error {
	if dto.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", model.ErrInvalidInput)
	}

	if dto.Nickname == "" {
		return fmt.Errorf("%w: nickname cannot be empty", model.ErrInvalidInput)
	}

	if dto.Likes < 0 || dto.Viewers < 0 {
		return fmt.Errorf("%w: likes and viewers cannot be negative", model.ErrInvalidInput)
	}

	if dto.Likes > dto.Viewers {
		return fmt.Errorf("%w: likes cannot be more then viewers", model.ErrInvalidInput)
	}

	user := model.NewUser(dto.Name, dto.Nickname, dto.Likes, dto.Viewers)

	if err := u.repo.Create(ctx, *user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (u *UserService) GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error) {
	if params.Sort != "" && params.Sort != "desc" && params.Sort != "asc" {
		return nil, -1, fmt.Errorf("%w: invalid sort parameter", model.ErrInvalidSort)
	}

	if params.Limit < 1 || params.Offset < 0 {
		return nil, -1, fmt.Errorf("%w: page or size cannot be negative or 0", model.ErrInvalidInput)
	}
	return u.repo.GetAll(ctx, params)
}

func (u *UserService) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	if nickname == "" {
		return nil, fmt.Errorf("%w: nickname cannot be empty", model.ErrInvalidInput)
	}

	return u.repo.GetUser(ctx, nickname)
}

func (u *UserService) ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error {
	if nickname == "" {
		return fmt.Errorf("%w: nickname cannot be empty", model.ErrInvalidInput)
	}

	if dto.Likes == nil && dto.Name == nil && dto.Nickname == nil && dto.Viewers == nil {
		return fmt.Errorf("%w: all data fields cannot be empty", model.ErrInvalidInput)
	}

	if dto.Likes != nil && *dto.Likes < 0 {
		return fmt.Errorf("%w: likes cannot be negative", model.ErrInvalidInput)
	}

	if dto.Name != nil && *dto.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", model.ErrInvalidInput)
	}

	if dto.Nickname != nil && *dto.Nickname == "" {
		return fmt.Errorf("%w: nickname cannot be empty", model.ErrInvalidInput)
	}

	if dto.Viewers != nil && *dto.Viewers < 0 {
		return fmt.Errorf("%w: viewers cannot be negative", model.ErrInvalidInput)
	}

	if dto.Likes != nil && dto.Viewers != nil && *dto.Likes > *dto.Viewers {
		return fmt.Errorf("%w: likes cannot be more than viewers", model.ErrInvalidInput)
	}

	if err := u.repo.ChangeData(ctx, nickname, dto); err != nil {
		return fmt.Errorf("failed to change data: %w", err)
	}

	return nil
}

func (u *UserService) Delete(ctx context.Context, nickname string) error {
	if nickname == "" {
		return fmt.Errorf("%w: nickname cannot be empty", model.ErrInvalidInput)
	}

	return u.repo.Delete(ctx, nickname)
}
