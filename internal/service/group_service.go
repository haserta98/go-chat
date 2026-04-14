package service

import (
	"errors"

	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
)

type GroupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
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
