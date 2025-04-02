package service

import (
	"context"
	"fmt"

	"github.com/DblMOKRQ/test_task/internal/adapters/agify"
	"github.com/DblMOKRQ/test_task/internal/adapters/genderize"
	"github.com/DblMOKRQ/test_task/internal/adapters/nationalize"
	"github.com/DblMOKRQ/test_task/internal/entity"
	"go.uber.org/zap"
)

type Repository interface {
	AddUser(ctx context.Context, u *entity.UserRequest) (int, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, u *entity.UserRequest, id int) error
	GetUserByID(ctx context.Context, id int) (*entity.FullUser, error)
	GetAllUsers(ctx context.Context) ([]*entity.FullUser, error)
	GetUsersByName(ctx context.Context, name string) ([]*entity.FullUser, error)
	GetUsersByAge(ctx context.Context, age int) ([]*entity.FullUser, error)
	GetUsersByGender(ctx context.Context, gender string) ([]*entity.FullUser, error)
	GetUsersByNationality(ctx context.Context, nationality string) ([]*entity.FullUser, error)
}

type Service struct {
	repo Repository
	log  *zap.Logger
	ctx  context.Context
}

func NewService(repo Repository, log *zap.Logger, ctx context.Context) *Service {
	return &Service{repo: repo, log: log, ctx: ctx}
}

func (s *Service) AddUser(user *entity.User) (int, error) {
	age, gender, national := s.getData(user.Name)
	userReq := entity.UserRequest{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: national,
	}

	id, err := s.repo.AddUser(s.ctx, &userReq)
	if err != nil {
		s.log.Error("failed to add user", zap.Error(err))
		return 0, err
	}
	s.log.Info("user added", zap.Int("id", id))
	return id, nil
}

func (s *Service) DeleteUser(id int) error {
	s.log.Info("user deleting", zap.Int("id", id))
	if err := s.repo.DeleteUser(s.ctx, id); err != nil {
		s.log.Error("failed to delete user", zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}
	s.log.Info("user deleted", zap.Int("id", id))
	return nil
}

func (s *Service) UpdateUser(newUser *entity.User, id int) error {
	age, gender, national := s.getData(newUser.Name)
	userReq := entity.UserRequest{
		Name:        newUser.Name,
		Surname:     newUser.Surname,
		Patronymic:  newUser.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: national,
	}
	s.log.Info("user updating", zap.Int("id", id))
	if err := s.repo.UpdateUser(s.ctx, &userReq, id); err != nil {
		s.log.Error("failed to update user", zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}
	s.log.Info("user updated", zap.Int("id", id))
	return nil
}

func (s *Service) GetUserByID(id int) (*entity.FullUser, error) {
	user, err := s.repo.GetUserByID(s.ctx, id)
	s.log.Info("user getting", zap.Int("id", id))
	if err != nil {
		s.log.Error("failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	s.log.Info("user get", zap.Int("id", id))
	return user, nil
}

func (s *Service) GetAllUsers() ([]*entity.FullUser, error) {
	s.log.Info("getting all users")
	users, err := s.repo.GetAllUsers(s.ctx)
	if err != nil {
		s.log.Error("failed to get all users", zap.Error(err))
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	s.log.Info("all users get", zap.Int("count", len(users)))
	return users, nil
}

func (s *Service) GetUsersByName(name string) ([]*entity.FullUser, error) {
	s.log.Info("getting users by name", zap.String("name", name))
	users, err := s.repo.GetUsersByName(s.ctx, name)
	if err != nil {
		s.log.Error("failed to get users by name", zap.Error(err))
		return nil, fmt.Errorf("failed to get users by name: %w", err)
	}
	s.log.Info("users by name get", zap.Int("count", len(users)))
	return users, nil
}

func (s *Service) GetUsersByAge(age int) ([]*entity.FullUser, error) {
	s.log.Info("getting users by age", zap.Int("age", age))
	users, err := s.repo.GetUsersByAge(s.ctx, age)
	if err != nil {
		s.log.Error("failed to get users by age", zap.Error(err))
		return nil, fmt.Errorf("failed to get users by age: %w", err)
	}
	s.log.Info("users by age get", zap.Int("count", len(users)))
	return users, nil
}

func (s *Service) GetUsersByGender(gender string) ([]*entity.FullUser, error) {
	s.log.Info("getting users by gender", zap.String("gender", gender))
	users, err := s.repo.GetUsersByGender(s.ctx, gender)
	if err != nil {
		s.log.Error("failed to get users by gender", zap.Error(err))
		return nil, fmt.Errorf("failed to get users by gender: %w", err)
	}
	s.log.Info("users by gender get", zap.Int("count", len(users)))
	return users, nil
}

func (s *Service) GetUsersByNationality(nationality string) ([]*entity.FullUser, error) {
	s.log.Info("getting users by nationality", zap.String("nationality", nationality))
	users, err := s.repo.GetUsersByNationality(s.ctx, nationality)
	if err != nil {
		s.log.Error("failed to get users by nationality", zap.Error(err))
		return nil, fmt.Errorf("failed to get users by nationality: %w", err)
	}
	s.log.Info("users by nationality get", zap.Int("count", len(users)))
	return users, nil
}

func (s *Service) getData(name string) (int, string, string) {
	age, err := agify.GetAge(name)
	if err != nil {
		s.log.Error("failed to get age", zap.Error(err))
		age = 0
	}
	gender, err := genderize.GetGender(name)
	if err != nil {
		s.log.Error("failed to get gender", zap.Error(err))
		gender = ""
	}
	national, err := nationalize.GetNationality(name)
	if err != nil {
		s.log.Error("failed to get nationality", zap.Error(err))
		national = ""
	}
	return age, gender, national
}
