package repository

import (
	"context"
	"testing"

	"go-tutorial/phase3/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "test@example.com",
		Role:     "user",
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		Username: "alice",
		Password: "pwd",
		Email:    "alice@example.com",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.GetByUsername(ctx, "alice")
	require.NoError(t, err)
	assert.Equal(t, "alice", found.Username)

	_, err = repo.GetByUsername(ctx, "notexist")
	assert.Error(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		user := &domain.User{
			Username: "user" + string(rune('0'+i)),
			Password: "pwd",
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	users, total, err := repo.List(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 5)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{Username: "todelete", Password: "pwd"}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
}
