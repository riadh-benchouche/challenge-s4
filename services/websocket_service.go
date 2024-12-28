package services

import (
	"backend/models"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	connection *websocket.Conn
	user       *models.User
}

var connections = make([]*WebSocketConnection, 0)

type WebSocketService struct {
	messageService *MessageService
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		messageService: NewMessageService(),
	}
}

func (s *WebSocketService) AcceptNewWebSocketConnection(ws *websocket.Conn, user *models.User) func(ws *websocket.Conn) {
	wsConnection := WebSocketConnection{
		connection: ws,
		user:       user,
	}
	connections = append(connections, &wsConnection)

	return func(ws *websocket.Conn) {
		_ = ws.Close()

		// Remove the closed WebSocket from the connection list
		for index, connection := range connections {
			if connection.connection == ws {
				connections = append(connections[:index], connections[index+1:]...)
				break
			}
		}
	}
}

func (s *WebSocketService) HandleWebSocketMessage(msg []byte, user *models.User, ws *websocket.Conn) error {
	var decodedMessage TypeMessage
	if err := json.Unmarshal(msg, &decodedMessage); err != nil {
		return err
	}

	switch decodedMessage.Type {
	case ClientBoundSendChatMessageType:
		if err := s.handleSendChatMessage(msg, user); err != nil {
			return err
		}
	case ClientBoundFetchChatMessageType:
		if err := s.handleFetchChatMessage(msg, ws); err != nil {
			return err
		}
	}

	return nil
}

func (s *WebSocketService) sendWebSocketMessage(data []byte, conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.TextMessage, data)
}

func (s *WebSocketService) BroadcastMessageToAssociation(associationID string) error {
	for _, connection := range connections {
		user := connection.user

		// Vérifiez si l'utilisateur est membre de l'association
		isMember := false
		for _, association := range user.Associations {
			if association.ID == associationID {
				isMember = true
				break
			}
		}

		if !isMember {
			continue
		}

		// Diffusez le message
		message := &models.Message{
			AssociationID: associationID,
			SenderID:      user.ID,
		}

		bytes, err := s.buildServerBoundSendChatMessage(message)
		if err != nil {
			return err
		}

		if err := s.sendWebSocketMessage(bytes, connection.connection); err != nil {
			return err
		}
	}

	return nil
}

func (s *WebSocketService) buildServerBoundSendChatMessage(message *models.Message) ([]byte, error) {
	response := ServerBoundSendChatMessage{
		TypeMessage: TypeMessage{
			Type: ServerBoundSendChatMessageType,
		},
		Content:       message.Content,
		Author:        &message.Sender,
		AssociationId: message.AssociationID,
		MessageId:     message.ID,
	}

	return json.Marshal(response)
}

func (s *WebSocketService) handleSendChatMessage(msg []byte, user *models.User) error {
	var receivedMessage ClientBoundSendChatMessage
	if err := json.Unmarshal(msg, &receivedMessage); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(receivedMessage); err != nil {
		return err
	}

	message := &models.Message{
		Content:       receivedMessage.Content,
		AssociationID: receivedMessage.AssociationId,
		SenderID:      user.ID,
	}

	if err := s.messageService.Create(message); err != nil {
		return err
	}

	return s.BroadcastMessageToAssociation(message.AssociationID)
}

func (s *WebSocketService) handleFetchChatMessage(msg []byte, ws *websocket.Conn) error {
	var receivedMessage ClientBoundFetchChatMessage
	if err := json.Unmarshal(msg, &receivedMessage); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(receivedMessage); err != nil {
		return err
	}

	messages, err := s.messageService.GetAll()
	if err != nil {
		return err
	}

	for _, message := range messages {
		if message.AssociationID != receivedMessage.AssociationId {
			continue
		}

		bytes, err := s.buildServerBoundSendChatMessage(&message)
		if err != nil {
			return err
		}

		if err := s.sendWebSocketMessage(bytes, ws); err != nil {
			return err
		}
	}

	return nil
}

type ClientBoundSendChatMessage struct {
	TypeMessage
	Content       string `json:"content" validate:"required,min=10,max=300"`
	AssociationId string `json:"association_id" validate:"required"`
}

type ServerBoundSendChatMessage struct {
	TypeMessage
	Content       string       `json:"content" validate:"required"`
	Author        *models.User `json:"author" validate:"required"`
	AssociationId string       `json:"association_id" validate:"required"`
	MessageId     string       `json:"message_id" validate:"required"`
}

type ClientBoundFetchChatMessage struct {
	TypeMessage
	AssociationId string `json:"association_id" validate:"required"`
}

type TypeMessage struct {
	Type string `json:"type" validate:"required"`
}

const (
	ClientBoundSendChatMessageType  string = "send_chat_message"
	ClientBoundFetchChatMessageType        = "fetch_chat_messages"
	ServerBoundSendChatMessageType         = "send_chat_message"
)
