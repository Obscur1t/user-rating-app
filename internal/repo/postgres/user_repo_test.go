package postgres

import (
	"context"
	"os"
	"path/filepath"
	"rating/internal/dto/request"
	"rating/internal/model"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	migrationPath := filepath.Join("..", "..", "..", "migrations", "20260225160021_init_users_table.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	require.NoError(t, err, "failed to read migration file")
	upSQL := strings.Split(string(migrationSQL), "-- +goose Down")[0]

	_, err = pool.Exec(ctx, string(upSQL))
	require.NoError(t, err, "failed to execute migration")

	return pool
}

func TestUserRepo_Create(t *testing.T) {
	pool := setupTestDB(t)
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	repo := NewUserRepo(pool)

	user := model.NewUser("name", "nickname", 50, 100)

	t.Run("success", func(t *testing.T) {
		err := repo.Create(ctx, *user)
		require.NoError(t, err)

		var testUser model.User
		err = pool.QueryRow(ctx, "SELECT name, nickname, likes, viewers FROM users WHERE nickname = $1", user.NickName).Scan(&testUser.Name, &testUser.NickName, &testUser.Likes, &testUser.Viewers)
		require.NoError(t, err)

		require.Equal(t, *user, testUser)
	})

	t.Run("already exists", func(t *testing.T) {
		err := repo.Create(ctx, *user)
		require.ErrorIs(t, err, model.ErrAlreadyExists)
	})
}

func TestUserRepo_GetAll(t *testing.T) {
	pool := setupTestDB(t)
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(func() { cancel() })
	repo := NewUserRepo(pool)

	expectedList := []model.User{
		{
			Id:       1,
			Name:     "name1",
			NickName: "nickname1",
			Likes:    1,
			Viewers:  11,
			Rating:   0.091,
		},
		{
			Id:       2,
			Name:     "name2",
			NickName: "nickname2",
			Likes:    2,
			Viewers:  22,
			Rating:   0.091,
		},
		{
			Id:       3,
			Name:     "name3",
			NickName: "nickname3",
			Likes:    2,
			Viewers:  22,
			Rating:   0.091,
		},
	}

	for _, u := range expectedList {
		err := repo.Create(ctx, u)
		require.NoError(t, err, "failed to create user")
	}

	t.Run("success", func(t *testing.T) {
		list, total, err := repo.GetAll(ctx, request.PaginationQuery{
			Sort:   "asc",
			Limit:  10,
			Offset: 0,
		})
		require.NoError(t, err)

		require.Equal(t, expectedList, list)
		require.Equal(t, 3, total)
	})

	t.Run("pagination limit", func(t *testing.T) {
		list, total, err := repo.GetAll(ctx, request.PaginationQuery{
			Sort:   "asc",
			Limit:  2,
			Offset: 0,
		})
		require.NoError(t, err)

		require.Equal(t, expectedList[:2], list)
		require.Equal(t, 3, total)
	})

	t.Run("pagination offset", func(t *testing.T) {
		list, total, err := repo.GetAll(ctx, request.PaginationQuery{
			Sort:   "asc",
			Limit:  2,
			Offset: 2,
		})
		require.NoError(t, err)

		require.Equal(t, expectedList[2:], list)
		require.Equal(t, 3, total)
	})
}

func TestUserRepo_GetUser(t *testing.T) {
	pool := setupTestDB(t)
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(func() { cancel() })
	repo := NewUserRepo(pool)
	expectedUser := model.User{
		Id:       1,
		Name:     "name",
		NickName: "nickname",
		Likes:    50,
		Viewers:  100,
		Rating:   0.5,
	}
	err := repo.Create(ctx, expectedUser)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		user, err := repo.GetUser(ctx, "nickname")
		require.NoError(t, err)

		require.Equal(t, expectedUser, *user)
	})

	t.Run("not found", func(t *testing.T) {
		user, err := repo.GetUser(ctx, "nil")
		require.ErrorIs(t, err, model.ErrNotFound)
		require.Nil(t, user)
	})
}

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }

func TestUserRepo_ChangeData(t *testing.T) {
	pool := setupTestDB(t)
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(func() { cancel() })
	repo := NewUserRepo(pool)

	user := model.NewUser("name", "nickname", 50, 100)
	err := repo.Create(ctx, *user)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		expectedUser := model.User{
			Id:       1,
			Name:     "name1",
			NickName: "nickname1",
			Likes:    100,
			Viewers:  200,
			Rating:   0.5,
		}

		err := repo.ChangeData(ctx, "nickname", request.UpdateUserDTO{
			Name:     ptrString("name1"),
			Nickname: ptrString("nickname1"),
			Likes:    ptrInt(100),
			Viewers:  ptrInt(200),
		})
		require.NoError(t, err)

		data, err := repo.GetUser(ctx, "nickname1")
		require.NoError(t, err)

		require.Equal(t, expectedUser, *data)
	})

	t.Run("likes more than viewers", func(t *testing.T) {
		err := repo.ChangeData(ctx, "nickname1", request.UpdateUserDTO{
			Name:     ptrString("name1"),
			Nickname: ptrString("nickname1"),
			Likes:    ptrInt(200),
			Viewers:  ptrInt(100),
		})
		require.ErrorIs(t, err, model.ErrInvalidInput)

	})

	t.Run("not found", func(t *testing.T) {
		err := repo.ChangeData(ctx, "nil", request.UpdateUserDTO{
			Name:     ptrString("name1"),
			Nickname: ptrString("nickname1"),
			Likes:    ptrInt(100),
			Viewers:  ptrInt(200),
		})
		require.ErrorIs(t, err, model.ErrNotFound)

	})
}

func TestUserRepo_Delete(t *testing.T) {
	pool := setupTestDB(t)
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(func() { cancel() })
	repo := NewUserRepo(pool)

	err := repo.Create(ctx, *model.NewUser("name", "nickname", 50, 100))
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := repo.Delete(ctx, "nickname")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(ctx, "nickname1")
		require.ErrorIs(t, err, model.ErrNotFound)
	})
}
