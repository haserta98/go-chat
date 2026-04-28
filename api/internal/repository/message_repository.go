package repository

import (
	"time"

	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/models"
)

type MessageRepository struct {
	db *cmd.DB
}

func NewMessageRepository(db *cmd.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.GetDB().Create(message).Error
}

func (r *MessageRepository) CreateGroupMessage(message *models.GroupMessage) error {
	return r.db.GetDB().Create(message).Error
}

func (r *MessageRepository) GetByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.GetDB().First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) GetAll() ([]models.Message, error) {
	var messages []models.Message
	err := r.db.GetDB().Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.GetDB().Updates(message).Error
}

func (r *MessageRepository) SeenBulk(messageIDs []string) error {
	return r.db.GetDB().Where("id IN ?", messageIDs).Update("seen_date", time.Now()).Error
}

func (r *MessageRepository) Delete(id string) (int64, error) {
	delete := r.db.GetDB().Delete(&models.Message{}, id)
	return delete.RowsAffected, delete.Error
}

func (r *MessageRepository) GetMessagesBetween(userID1, userID2 string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.GetDB().
		Raw("SELECT * FROM messages WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY created_at ASC", userID1, userID2, userID2, userID1).
		Scan(&messages).Error
	return messages, err
}

func (r *MessageRepository) UpsertContact(userID, contactID string) error {
	now := time.Now()

	// Create/Update A -> B
	err := r.db.GetDB().Save(&models.UserContact{
		UserID:        userID,
		ContactID:     contactID,
		LastMessageAt: now,
	}).Error
	if err != nil {
		return err
	}

	// Create/Update B -> A
	return r.db.GetDB().Save(&models.UserContact{
		UserID:        contactID,
		ContactID:     userID,
		LastMessageAt: now,
	}).Error
}

func (r *MessageRepository) GetGroupMessages(groupID string) ([]models.GroupMessage, error) {
	var messages []models.GroupMessage
	err := r.db.GetDB().Raw("SELECT * FROM group_messages WHERE to_id = ? ORDER BY created_at ASC", groupID).Scan(&messages).Error
	return messages, err
}
