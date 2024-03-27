package repository

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	"github.com/ovrrtd/openidea-bank/internal/model/entity"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	Register(ctx context.Context, user entity.User) (*entity.User, int, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, int, error)
	FindByID(ctx context.Context, id int64) (*entity.User, int, error)
	UpdateByID(ctx context.Context, user entity.User) (*entity.User, int, error)
}

func NewUserRepository(logger zerolog.Logger, db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type UserRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *UserRepositoryImpl) Register(ctx context.Context, newUser entity.User) (*entity.User, int, error) {
	newUser.CreatedAt = time.Now().UnixMilli()
	newUser.UpdatedAt = time.Now().UnixMilli()

	err := r.db.QueryRowContext(ctx, "INSERT INTO users (email, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newUser.Email, newUser.Password, newUser.Name, newUser.CreatedAt, newUser.UpdatedAt).Scan(&newUser.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &newUser, http.StatusCreated, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, int, error) {
	var user entity.User

	row := r.db.QueryRowContext(ctx, "SELECT id, email, password, name,  created_at, updated_at FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return &user, http.StatusOK, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.User, int, error) {
	var user entity.User

	row := r.db.QueryRowContext(ctx, "SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return &user, http.StatusOK, nil
}

func (r *UserRepositoryImpl) UpdateByID(ctx context.Context, user entity.User) (*entity.User, int, error) {
	user.UpdatedAt = time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET 
			email = $1, password = $2, name = $3, updated_at = $4 
			WHERE id = $5
	`, user.Email, user.Password, user.Name, user.UpdatedAt, user.ID)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return &user, http.StatusOK, nil
}
