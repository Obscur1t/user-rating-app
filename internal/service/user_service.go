package service

import (
	"context"
	"fmt"
	"rating/internal/dto/request"
	"rating/internal/model"
)

type UserStore interface {
	Create(ctx context.Context, user model.User) error
	GetAll(ctx context.Context, sort string) ([]model.User, error)
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
		return fmt.Errorf("name cannot be empty")
	}

	if dto.Nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if dto.Likes < 0 || dto.Viewers < 0 {
		return fmt.Errorf("likes and viewers cannot be negative")
	}

	if dto.Likes > dto.Viewers {
		return fmt.Errorf("likes cannot be more then viewers")
	}

	user := model.NewUser(dto.Name, dto.Nickname, dto.Likes, dto.Viewers)

	if err := u.repo.Create(ctx, *user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (u *UserService) GetAll(ctx context.Context, sort string) ([]model.User, error) {
	return u.repo.GetAll(ctx, sort)
}

func (u *UserService) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	if nickname == "" {
		return nil, fmt.Errorf("nickname cannot be empty")
	}

	user, err := u.repo.GetUser(ctx, nickname)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (u *UserService) ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if dto.Likes == nil && dto.Name == nil && dto.Nickname == nil && dto.Viewers == nil {
		return fmt.Errorf("all data fields cannot be empty")
	}

	if dto.Likes != nil && *dto.Likes < 0 {
		return fmt.Errorf("likes cannot be negative")
	}

	if dto.Name != nil && *dto.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if dto.Nickname != nil && *dto.Nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if dto.Viewers != nil && *dto.Viewers < 0 {
		return fmt.Errorf("viewers cannot be negative")
	}

	if err := u.repo.ChangeData(ctx, nickname, dto); err != nil {
		return fmt.Errorf("failed to change data: %w", err)
	}

	return nil
}

func (u *UserService) Delete(ctx context.Context, nickname string) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	return u.repo.Delete(ctx, nickname)
}
