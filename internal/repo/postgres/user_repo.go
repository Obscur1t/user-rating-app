package postgres

import (
	"context"
	"errors"
	"fmt"
	"rating/internal/dto/request"
	"rating/internal/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *UserRepo) scanUser(rows pgx.Rows) ([]model.User, error) {
	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to scan user data", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user model.User) error {
	query := "INSERT INTO users (name, nickname, likes, viewers) VALUES($1, $2, $3, $4)"

	_, err := r.pool.Exec(ctx, query, user.Name, user.NickName, user.Likes, user.Viewers)

	var pgxErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return fmt.Errorf("%w: nickname %s already exists", model.ErrAlreadyExists, user.NickName)
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetAll(ctx context.Context, params request.PaginationQuery) ([]model.User, int, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users"

	var totalCount int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&totalCount); err != nil {
		return nil, -1, fmt.Errorf("failed to get total count users: %w", err)
	}

	if params.Sort == "desc" {
		query += " ORDER BY rating DESC"
	} else if params.Sort == "asc" {
		query += " ORDER BY rating ASC"
	}

	query += " LIMIT $1 OFFSET $2"

	rows, err := r.pool.Query(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	userList, err := r.scanUser(rows)

	return userList, totalCount, err
}

func (r *UserRepo) GetUser(ctx context.Context, nickname string) (*model.User, error) {
	query := "SELECT id, name, nickname, likes, viewers, rating FROM users WHERE nickname = $1"

	var user model.User
	err := r.pool.QueryRow(ctx, query, nickname).Scan(&user.Id, &user.Name, &user.NickName, &user.Likes, &user.Viewers, &user.Rating)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found", model.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) ChangeData(ctx context.Context, nickname string, dto request.UpdateUserDTO) error {
	var args []any
	var sets []string

	if dto.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *dto.Name)
	}

	if dto.Likes != nil {
		sets = append(sets, fmt.Sprintf("likes = $%d", len(args)+1))
		args = append(args, *dto.Likes)
	}

	if dto.Viewers != nil {
		sets = append(sets, fmt.Sprintf("viewers = $%d", len(args)+1))
		args = append(args, *dto.Viewers)
	}

	if dto.Nickname != nil {
		sets = append(sets, fmt.Sprintf("nickname = $%d", len(args)+1))
		args = append(args, *dto.Nickname)
	}

	args = append(args, nickname)

	query := fmt.Sprintf("UPDATE users SET %s WHERE %s ", strings.Join(sets, ", "), fmt.Sprintf("nickname = $%d", len(args)))

	cmdTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23514" {
				return fmt.Errorf("%w: likes can't be more than viewers", model.ErrInvalidInput)
			}
		}
		return fmt.Errorf("failed to update data: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: user not found", model.ErrNotFound)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, nickname string) error {
	query := "DELETE FROM users WHERE nickname = $1"

	cmdTag, err := r.pool.Exec(ctx, query, nickname)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: user not found", model.ErrNotFound)
	}

	return nil

}
