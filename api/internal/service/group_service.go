package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/ws"
)

type JoinGroup struct {
	GroupID string `json:"groupId"`
}

type LeaveGroup struct {
	GroupID string `json:"groupId"`
}

type GroupService struct {
	repo      *repository.GroupRepository
	wsGateway *ws.WsGateway
}

func NewGroupService(repo *repository.GroupRepository, wsGateway *ws.WsGateway) *GroupService {
	return &GroupService{repo: repo, wsGateway: wsGateway}
}

func (s *GroupService) Start() {

}

func (s *GroupService) CreateGroup(name string) (*models.Group, error) {
	if name == "" {
		return nil, errors.New("name alanı zorunludur")
	}
	group := &models.Group{
		ID:   uuid.New().String(),
		Name: name,
	}
	return group, s.repo.Create(group)
}

func (s *GroupService) AddMember(groupID, userID string) error {
	return s.repo.AddMember(&models.GroupMember{
		ID:      uuid.New().String(),
		GroupID: groupID,
		UserID:  userID,
	})
}

func (s *GroupService) RemoveMember(groupID, userID string) error {
	return s.repo.RemoveMember(&models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
	})
}

func (s *GroupService) GetGroupByID(id string) (*models.Group, error) {
	if id == "" {
		return nil, errors.New("id alanı zorunludur")
	}
	return s.repo.GetByID(id)
}

func (s *GroupService) GetAllGroups() ([]models.Group, error) {
	return s.repo.GetAll()
}

func (s *GroupService) UpdateGroup(id, name string) error {
	group, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	group.Name = name
	return s.repo.Update(group)
}

func (s *GroupService) DeleteGroup(id string) error {
	if id == "" {
		return errors.New("id alanı zorunludur")
	}
	count, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("grup bulunamadı")
	}
	return nil
}

func (s *GroupService) GetMyGroups(userID string) ([]models.Group, error) {
	return s.repo.GetMyGroups(userID)
}
