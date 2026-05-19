package repository

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"

	"weather-api/internal/model"
)

type SQLiteCityRepository struct {
	db *sql.DB
}

func NewSQLiteCityRepository(db *sql.DB) *SQLiteCityRepository {
	return &SQLiteCityRepository{db: db}
}

func (r *SQLiteCityRepository) Create(ctx context.Context, userID int64, name string) (*model.City, error) {
	now := time.Now().UTC()

	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO cities (user_id, name, created_at) VALUES (?, ?, ?)`,
		userID,
		name,
		now,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.City{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
	}, nil
}

func (r *SQLiteCityRepository) ListByUserID(ctx context.Context, userID int64) ([]*model.City, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, created_at FROM cities WHERE user_id = ? ORDER BY id ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []*model.City
	for rows.Next() {
		city := &model.City{}
		if err := rows.Scan(&city.ID, &city.UserID, &city.Name, &city.CreatedAt); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, rows.Err()
}

func (r *SQLiteCityRepository) GetByID(ctx context.Context, id int64) (*model.City, error) {
	city := &model.City{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, created_at FROM cities WHERE id = ?`,
		id,
	).Scan(&city.ID, &city.UserID, &city.Name, &city.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return city, nil
}

func (r *SQLiteCityRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cities WHERE id = ?`, id)
	return err
}

func (r *SQLiteCityRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cities WHERE user_id = ?`, userID)
	return err
}
