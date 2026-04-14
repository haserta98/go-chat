package service

import (
	"errors"

	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
)

type UserService interface {
	CreateUser(name, email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id, name, email string) error
	DeleteUser(id string) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(name, email string) (*models.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name ve Email alanları zorunludur")
	}
	user := &models.User{
		Name:  name,
		Email: email,
	}
	return user, s.repo.Create(user)
}
func (s *UserServiceImpl) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}
func (s *UserServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}
func (s *UserServiceImpl) UpdateUser(id, name, email string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	return s.repo.Update(user)
}
func (s *UserServiceImpl) DeleteUser(id string) error {
	count, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("kullanıcı bulunamadı")
	}
	return nil
}
