package repo

import (
	"context"
	"rating/internal/model"
)

type RepoInterface interface {
	Create(ctx context.Context, user model.User) error
	GetAll(ctx context.Context) ([]model.User, error)
	GetFilteredByRatingDESC(ctx context.Context) ([]model.User, error)
	GetFilteredByRatingASC(ctx context.Context) ([]model.User, error)
	GetUser(ctx context.Context, nickname string) (*model.User, error)
	ChangeLikes(ctx context.Context, nickname string, value int) error
	ChangeViewers(ctx context.Context, nickname string, value int) error
	ChangeName(ctx context.Context, nickname string, value string) error
	ChangeNickname(ctx context.Context, nickname string, value string) error
	Delete(ctx context.Context, nickname string) error
}
