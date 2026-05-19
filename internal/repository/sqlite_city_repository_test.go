package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteCityRepositoryCreateAndRead(t *testing.T) {
	db := newTestSQLiteDB(t)
	repo := NewSQLiteCityRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, 42, "Almaty")
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	stored, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, created.ID, stored.ID)
	assert.Equal(t, int64(42), stored.UserID)
	assert.Equal(t, "Almaty", stored.Name)

	cities, err := repo.ListByUserID(ctx, 42)
	require.NoError(t, err)
	require.Len(t, cities, 1)
	assert.Equal(t, created.ID, cities[0].ID)
	assert.Equal(t, "Almaty", cities[0].Name)

	require.NoError(t, repo.Delete(ctx, created.ID))

	stored, err = repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Nil(t, stored)

	cities, err = repo.ListByUserID(ctx, 42)
	require.NoError(t, err)
	assert.Empty(t, cities)
}

func newTestSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	_, err = db.Exec(`
		CREATE TABLE cities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}
