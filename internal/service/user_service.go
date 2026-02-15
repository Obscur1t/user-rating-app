package service

import (
	"context"
	"fmt"
	"rating/internal/dto/request"
	"rating/internal/model"
	"rating/internal/repo"
)

type UserService struct {
	repo repo.RepoInterface
}

func NewUserService(repo repo.RepoInterface) *UserService {
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
		return fmt.Errorf("failed to create use %v", err)
	}

	return nil
}

func (u *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *UserService) GetFilteredByRatingDESC(ctx context.Context) ([]model.User, error) {
	return u.repo.GetFilteredByRatingDESC(ctx)
}

func (u *UserService) GetFilteredByRatingASC(ctx context.Context) ([]model.User, error) {
	return u.repo.GetFilteredByRatingASC(ctx)
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

func (u *UserService) ChangeLikes(ctx context.Context, nickname string, value int) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if value < 0 {
		return fmt.Errorf("likes cannot be negative")
	}

	user, err := u.repo.GetUser(ctx, nickname)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if value > user.Viewers {
		return fmt.Errorf("likes cannot be more than viewers")
	}

	if err := u.repo.ChangeLikes(ctx, nickname, value); err != nil {
		return fmt.Errorf("failed to change likes %v", err)
	}

	return nil
}

func (u *UserService) ChangeViewers(ctx context.Context, nickname string, value int) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if value < 0 {
		return fmt.Errorf("viewers cannot be negative")
	}

	user, err := u.repo.GetUser(ctx, nickname)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if value < user.Viewers {
		return fmt.Errorf("viewers cannot be less than previous value")
	}

	if value < user.Likes {
		return fmt.Errorf("viewers cannot be less than likes")
	}

	if err := u.repo.ChangeViewers(ctx, nickname, value); err != nil {
		return fmt.Errorf("failed to change viewers %v", err)
	}

	return nil
}

func (u *UserService) ChangeName(ctx context.Context, nickname, value string) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if value == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if err := u.repo.ChangeName(ctx, nickname, value); err != nil {
		return fmt.Errorf("failed to change name %v", err)
	}

	return nil
}

func (u *UserService) ChangeNickname(ctx context.Context, nickname, value string) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if value == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	if err := u.repo.ChangeNickname(ctx, nickname, value); err != nil {
		return fmt.Errorf("failed to change nickname %v", err)
	}

	return nil
}

func (u *UserService) Delete(ctx context.Context, nickname string) error {
	if nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}

	return u.repo.Delete(ctx, nickname)
}
