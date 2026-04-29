package repository

import (
	"errors"

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
	err := r.db.GetDB().First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByName(name string) (*models.User, error) {
	var user models.User
	err := r.db.GetDB().Where("name = ?", name).First(&user).Error
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

func (r *UserRepository) GetContacts(userID string) ([]models.User, error) {
	var users []models.User
	err := r.db.GetDB().Raw(`
		SELECT u.* FROM users u 
		JOIN user_contacts uc ON uc.contact_id = u.id
		WHERE uc.user_id = ?
		ORDER BY uc.last_message_at DESC
	`, userID).Scan(&users).Error
	return users, err
}

func (r *UserRepository) AddContact(userID, contactID string) error {
	var existing models.UserContact
	if err := r.db.GetDB().Where("user_id = ? AND contact_id = ?", userID, contactID).First(&existing).Error; err != nil {
		contact := &models.UserContact{
			UserID:    userID,
			ContactID: contactID,
		}
		return r.db.GetDB().Create(contact).Error
	}
	return errors.New("contact already exists")
}

func (r *UserRepository) RemoveContact(userID, contactID string) error {
	return r.db.GetDB().Delete(&models.UserContact{}, "user_id = ? AND contact_id = ?", userID, contactID).Error
}

// GetContactOwners returns user IDs who have the given contactID in their contacts list.
// This is a reverse lookup used for presence notifications.
func (r *UserRepository) GetContactOwners(contactID string) ([]string, error) {
	var userIDs []string
	err := r.db.GetDB().Model(&models.UserContact{}).
		Where("contact_id = ?", contactID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

