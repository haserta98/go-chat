package repository

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/models"
)

type GroupRepository struct {
	db *cmd.DB
}

func NewGroupRepository(db *cmd.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *models.Group) error {
	return r.db.GetDB().Create(group).Error
}

func (r *GroupRepository) GetByID(id string) (*models.Group, error) {
	var group models.Group
	err := r.db.GetDB().First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) GetMembers(groupID string) ([]models.User, error) {
	var members []models.User
	err := r.db.GetDB().Where("group_id = ?", groupID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *GroupRepository) AddMember(member *models.GroupMember) error {
	return r.db.GetDB().Create(member).Error
}

func (r *GroupRepository) RemoveMember(member *models.GroupMember) error {
	return r.db.GetDB().Delete(member).Error
}

func (r *GroupRepository) GetAll() ([]models.Group, error) {
	var groups []models.Group
	err := r.db.GetDB().Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GroupRepository) Update(group *models.Group) error {
	return r.db.GetDB().Updates(group).Error
}

func (r *GroupRepository) Delete(id string) (int64, error) {
	delete := r.db.GetDB().Delete(&models.Group{}, id)
	return delete.RowsAffected, delete.Error
}

func (r *GroupRepository) GetMyGroups(userID string) ([]models.Group, error) {
	var groups []models.Group
	err := r.db.GetDB().Raw(`
		SELECT g.* FROM groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE gm.user_id = ?
	`, userID).Scan(&groups).Error
	return groups, err
}
