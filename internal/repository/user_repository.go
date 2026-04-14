package repository

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/models"
)

type UserRepository struct {
	db *cmd.DB
}

func NewUserRepository(db *cmd.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.GetDB().Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.GetDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.GetDB().Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.GetDB().Updates(user).Error
}

func (r *UserRepository) Delete(id string) (int64, error) {
	delete := r.db.GetDB().Delete(&models.User{}, id)
	return delete.RowsAffected, delete.Error
}
