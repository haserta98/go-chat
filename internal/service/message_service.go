package service

import (
	"encoding/json"
	"log"

	"github.com/bytedance/sonic"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/ws"
)

type SendMessage struct {
	To      string          `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

type SendGroupMessage struct {
	GroupID string          `json:"groupId"`
	Payload json.RawMessage `json:"payload"`
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

	s.wsGateway.Manager.RegisterEventHandler("send_message", func(client *ws.WsClient, payload json.RawMessage) {
		var sendMessage SendMessage
		if err := sonic.Unmarshal(payload, &sendMessage); err != nil {
			log.Printf("Invalid send message payload: %v", err)
			return
		}

		s.wsGateway.Manager.SendSmart(sendMessage.To, sendMessage.Payload)

		message := s.BuildMessage(string(sendMessage.Payload), client.UserID, sendMessage.To)
		s.PersistMessage(message)
	})

	s.wsGateway.Manager.RegisterEventHandler("send_group_message", func(client *ws.WsClient, payload json.RawMessage) {
		var sendGroupMessage SendGroupMessage
		if err := sonic.Unmarshal(payload, &sendGroupMessage); err != nil {
			log.Printf("Invalid send group message payload: %v", err)
			return
		}
		go s.wsGateway.Manager.BroadcastToGroup(client.UserID, sendGroupMessage.GroupID, sendGroupMessage.Payload)

		message := s.BuildGroupMessage(string(sendGroupMessage.Payload), client.UserID, sendGroupMessage.GroupID)
		s.PersistGroupMessage(message)
	})
}

func (s *MessageService) SeenBulk(messageIDs []string) error {
	return s.messageRepository.SeenBulk(messageIDs)
}

func (s *MessageService) PersistMessage(message *models.Message) {
	if err := s.messageRepository.Create(message); err != nil {
		log.Printf("Failed to persist message: %v", err)
	}
}

func (s *MessageService) PersistGroupMessage(message *models.GroupMessage) {
	if err := s.messageRepository.CreateGroupMessage(message); err != nil {
		log.Printf("Failed to persist group message: %v", err)
	}
}

func (s *MessageService) BuildMessage(content string, fromID string, toUserID string) *models.Message {
	return &models.Message{
		Content:  content,
		FromID:   fromID,
		ToUserID: toUserID,
	}
}

func (s *MessageService) BuildGroupMessage(content string, fromID string, toGroupID string) *models.GroupMessage {
	return &models.GroupMessage{
		Content:   content,
		FromID:    fromID,
		ToGroupID: toGroupID,
	}
}
