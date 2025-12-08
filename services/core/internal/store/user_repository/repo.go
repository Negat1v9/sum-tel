package user_repository

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	row := r.db.QueryRowxContext(
		ctx,
		createUserQuery,
		user.TelegramID,
		user.Username,
		user.IsActive,
		user.Role,
	)

	if err := row.StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.GetContext(ctx, user, getUserByIDQuery, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]model.User, error) {
	users := []model.User{}
	err := r.db.SelectContext(ctx, &users, getAllUsersQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	row := r.db.QueryRowxContext(
		ctx,
		updateUserQuery,
		user.ID,
		user.TelegramID,
		user.Username,
		user.IsActive,
		user.Role,
	)

	if err := row.StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.GetContext(ctx, user, deleteUserQuery, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
