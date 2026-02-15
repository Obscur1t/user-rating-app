package usersrepo

import (
	"context"
	"fmt"
	"rating/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) Create(ctx context.Context, user model.User) error {
	query := "INSERT INTO users (name, nickname, likes, viewers) VALUES($1, $2, $3, $4)"

	_, err := r.pool.Exec(ctx, query, user.Name, user.NickName, user.Likes, user.Viewers)
	if err != nil {
		return fmt.Errorf("failed to create user %v", err)
	}

	return nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users %v", err)
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user data %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) GetFilteredByRatingDESC(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users ORDER BY rating DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users %v", err)
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user data %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) GetFilteredByRatingASC(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users ORDER BY rating ASC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users %v", err)
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user data %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users WHERE nickname = $1 "

	var user model.User
	err := r.pool.QueryRow(ctx, query, nickname).Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %v", err)
	}

	return &user, nil
}

func (r *UserRepo) ChangeLikes(ctx context.Context, nickname string, value int) error {
	query := "UPDATE users SET likes = $1 WHERE nickname = $2 AND $1 <= viewers "

	cmdTag, err := r.pool.Exec(ctx, query, value, nickname)
	if err != nil {
		return fmt.Errorf("failed to update likes")
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user cannot be updated")
	}

	return nil
}

func (r *UserRepo) ChangeViewers(ctx context.Context, nickname string, value int) error {
	query := "UPDATE users SET viewers = $1 WHERE nickname = $2 AND $1 >= viewers AND $1 >= likes"

	cmdTag, err := r.pool.Exec(ctx, query, value, nickname)
	if err != nil {
		return fmt.Errorf("failed to update viewers")
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user cannot be updated")
	}

	return nil
}

func (r *UserRepo) ChangeName(ctx context.Context, nickname string, value string) error {
	query := "UPDATE users SET name = $1 WHERE nickname = $2 "

	cmdTag, err := r.pool.Exec(ctx, query, value, nickname)
	if err != nil {
		return fmt.Errorf("failed to update name %v", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found %v", err)
	}

	return nil
}

func (r *UserRepo) ChangeNickname(ctx context.Context, nickname string, value string) error {
	query := "UPDATE users SET nickname = $1 WHERE nickname = $2 "

	cmdTag, err := r.pool.Exec(ctx, query, value, nickname)
	if err != nil {
		return fmt.Errorf("failed to update nickname %v", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found %v", err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, nickname string) error {
	query := "DELETE FROM users WHERE nickname = $1"

	cmdTag, err := r.pool.Exec(ctx, query, nickname)
	if err != nil {
		return fmt.Errorf("failed to delete user %v", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
