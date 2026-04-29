package service

import (
	"encoding/json"
	"log"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/ws"
)

type SendMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type SendGroupMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type MessageService struct {
	wsGateway         *ws.WsGateway
	messageRepository *repository.MessageRepository
}

func NewMessageService(wsGateway *ws.WsGateway, messageRepository *repository.MessageRepository) *MessageService {
	return &MessageService{
		wsGateway:         wsGateway,
		messageRepository: messageRepository,
	}
}

func (s *MessageService) RegisterEventHandlers() {
	s.wsGateway.Manager.RegisterEventHandler("send_message", s.sendMessage)
	s.wsGateway.Manager.RegisterEventHandler("send_group_message", s.sendGroupMessage)
}

func (s *MessageService) sendGroupMessage(client *ws.WsClient, payload json.RawMessage) {
	var sendGroupMessage SendGroupMessage
	if err := sonic.Unmarshal(payload, &sendGroupMessage); err != nil {
		log.Printf("Invalid send group message payload: %v", err)
		return
	}

	go s.wsGateway.Manager.SendSmartGroup(client, string(sendGroupMessage.To), payload)

	message := s.BuildGroupMessage(sendGroupMessage, client)
	err := s.PersistGroupMessage(message)
	if err != nil {
		log.Printf("Failed to persist group message: %v", err)
	}
}

func (s *MessageService) sendMessage(client *ws.WsClient, payload json.RawMessage) {
	var sendMessage SendMessage
	if err := sonic.Unmarshal(payload, &sendMessage); err != nil {
		log.Printf("Invalid send message payload: %v", err)
		return
	}
	go s.wsGateway.Manager.SendSmart(sendMessage.To, payload)
	go s.PersistMessage(s.BuildMessage(sendMessage, client))
}

func (s *MessageService) SeenBulk(messageIDs []string) error {
	return s.messageRepository.SeenBulk(messageIDs)
}

func (s *MessageService) PersistMessage(message *models.Message) {
	if err := s.messageRepository.Create(message); err != nil {
		log.Printf("Failed to persist message: %v", err)
	} else {
		if err := s.messageRepository.UpsertContact(message.FromID, message.ToID); err != nil {
			log.Printf("Failed to upsert contact metadata: %v", err)
		}
	}
}

func (s *MessageService) PersistGroupMessage(message *models.GroupMessage) error {
	if err := s.messageRepository.CreateGroupMessage(message); err != nil {
		log.Printf("Failed to persist group message: %v", err)
		return err
	}
	return nil
}

func (s *MessageService) BuildMessage(sendMessage SendMessage, fromUser *ws.WsClient) *models.Message {
	return &models.Message{
		ID:      uuid.New().String(),
		Content: sendMessage.Content,
		FromID:  fromUser.UserID,
		ToID:    sendMessage.To,
	}
}

func (s *MessageService) BuildGroupMessage(sendGroupMessage SendGroupMessage, fromUser *ws.WsClient) *models.GroupMessage {
	return &models.GroupMessage{
		ID:      uuid.New().String(),
		Content: sendGroupMessage.Content,
		FromID:  fromUser.UserID,
		ToID:    sendGroupMessage.To,
	}
}

func (s *MessageService) GetMessagesBetween(userID1, userID2 string) ([]models.Message, error) {
	return s.messageRepository.GetMessagesBetween(userID1, userID2)
}

func (s *MessageService) GetGroupMessages(groupID string) ([]models.GroupMessage, error) {
	return s.messageRepository.GetGroupMessages(groupID)
}
