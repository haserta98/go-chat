package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/bytedance/sonic"
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
	repo      repository.GroupRepository
	wsGateway *ws.WsGateway
}

func NewGroupService(repo repository.GroupRepository, wsGateway *ws.WsGateway) *GroupService {
	return &GroupService{repo: repo, wsGateway: wsGateway}
}

func (s *GroupService) Start() {
	s.wsGateway.Manager.RegisterEventHandler("leave_group", func(client *ws.WsClient, payload json.RawMessage) {
		var leaveGroup LeaveGroup
		if err := sonic.Unmarshal(payload, &leaveGroup); err != nil {
			log.Printf("Invalid leave group payload: %v", err)
			return
		}
		s.wsGateway.Manager.LeaveGroup(leaveGroup.GroupID, client)
	})

	s.wsGateway.Manager.RegisterEventHandler("join_group", func(client *ws.WsClient, payload json.RawMessage) {
		var joinGroup JoinGroup
		if err := sonic.Unmarshal(payload, &joinGroup); err != nil {
			log.Printf("Invalid join group payload: %v", err)
			return
		}
		s.wsGateway.Manager.JoinGroup(joinGroup.GroupID, client)
	})
}

func (s *GroupService) CreateGroup(name string) (*models.Group, error) {
	if name == "" {
		return nil, errors.New("name alanı zorunludur")
	}
	group := &models.Group{
		Name: name,
	}
	return group, s.repo.Create(group)
}

func (s *GroupService) AddMember(groupID, userID string) error {
	return s.repo.AddMember(&models.GroupMember{
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
