package repository

import (
	"context"
	"fmt"

	"github.com/DblMOKRQ/test_task/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	storage *pgxpool.Pool
}

func NewRepository(storage *pgxpool.Pool) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) AddUser(ctx context.Context, u *entity.UserRequest) (int, error) {
	lastInsertId := 0
	err := r.storage.QueryRow(ctx, "INSERT INTO users (name, surname, patronymic,age,gender,nationality) VALUES ($1, $2, $3,$4,$5,$6) RETURNING id", u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.Nationality).Scan(&lastInsertId)
	if err != nil {
		return 0, fmt.Errorf("failed to add user: %w", err)
	}

	return lastInsertId, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.storage.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, u *entity.UserRequest, id int) error {
	_, err := r.storage.Exec(ctx, "UPDATE users SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE id = $7", u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.Nationality, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *Repository) GetAllUsers(ctx context.Context) ([]*entity.FullUser, error) {
	var users []*entity.FullUser
	rows, err := r.storage.Query(ctx, "SELECT id, name, surname, patronymic, age, gender, nationality FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.FullUser
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*entity.FullUser, error) {
	users, err := r.getUsersBy(ctx, "id", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user.id: %w", err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (r *Repository) GetUsersByName(ctx context.Context, name string) ([]*entity.FullUser, error) {
	users, err := r.getUsersBy(ctx, "name", name)
	if err != nil {
		return nil, fmt.Errorf("failed to get users.name: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUsersByAge(ctx context.Context, age int) ([]*entity.FullUser, error) {
	users, err := r.getUsersBy(ctx, "age", age)
	if err != nil {
		return nil, fmt.Errorf("failed to get users.age: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUsersByGender(ctx context.Context, gender string) ([]*entity.FullUser, error) {
	users, err := r.getUsersBy(ctx, "gender", gender)
	if err != nil {
		return nil, fmt.Errorf("failed to get users.gender: %w", err)
	}
	return users, nil
}

func (r *Repository) GetUsersByNationality(ctx context.Context, nationality string) ([]*entity.FullUser, error) {
	users, err := r.getUsersBy(ctx, "nationality", nationality)
	if err != nil {
		return nil, fmt.Errorf("failed to get users.Nationality: %w", err)
	}
	return users, nil
}

func (r *Repository) getUsersBy(ctx context.Context, where string, value any) ([]*entity.FullUser, error) {
	var users []*entity.FullUser
	rows, err := r.storage.Query(ctx, "SELECT id, name, surname, patronymic, age, gender, nationality FROM users WHERE "+where+" = $1", value)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.FullUser
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}
