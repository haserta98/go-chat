package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return []byte(s)
}

type UserService interface {
	CreateUser(name, email, password string) (*models.User, error)
	LoginUser(name, password string) (string, *models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id, name, email string) error
	DeleteUser(id string) error
	GetContacts(userID string) ([]models.User, error)
	AddContact(userID, contactID string) error
	RemoveContact(userID, contactID string) error
}

type UserServiceImpl struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(name, email, password string) (*models.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email ve password alanları zorunludur")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("şifre hashlenirken hata oluştu")
	}

	user := &models.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	return user, s.repo.Create(user)
}

func (s *UserServiceImpl) LoginUser(name, password string) (string, *models.User, error) {
	user, err := s.repo.GetByName(name)
	if err != nil {
		return "", nil, errors.New("kullanıcı bulunamadı")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("geçersiz şifre")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret())
	if err != nil {
		return "", nil, errors.New("token oluşturulamadı")
	}

	return tokenString, user, nil
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

func (s *UserServiceImpl) AddContact(userID, contactID string) error {
	if userID == contactID {
		return errors.New("kendini kişi olarak ekleyemezsin")
	}
	return s.repo.AddContact(userID, contactID)
}

func (s *UserServiceImpl) RemoveContact(userID, contactID string) error {
	return s.repo.RemoveContact(userID, contactID)
}

func (s *UserServiceImpl) GetContacts(userID string) ([]models.User, error) {
	return s.repo.GetContacts(userID)
}
