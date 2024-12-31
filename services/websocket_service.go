package services

import (
	"backend/database"
	"backend/enums"
	"backend/models"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	connection *websocket.Conn
	user       *models.User
}

// Liste des connexions actives
var connections = make([]*WebSocketConnection, 0)

type WebSocketService struct {
	messageService     *MessageService
	associationService *AssociationService
}

// Initialisation du service WebSocket
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		messageService:     NewMessageService(),
		associationService: NewAssociationService(),
	}
}

// Accepter une nouvelle connexion WebSocket
func (s *WebSocketService) AcceptNewWebSocketConnection(ws *websocket.Conn, user *models.User) func(ws *websocket.Conn) {
	wsConnection := WebSocketConnection{
		connection: ws,
		user:       user,
	}
	connections = append(connections, &wsConnection)

	return func(ws *websocket.Conn) {
		_ = ws.Close()
		// Supprimer la connexion fermée de la liste
		for index, connection := range connections {
			if connection.connection == ws {
				connections = append(connections[:index], connections[index+1:]...)
				break
			}
		}
	}
}

// Gérer un message reçu via WebSocket
func (s *WebSocketService) HandleWebSocketMessage(msg []byte, user *models.User, ws *websocket.Conn) error {
	var receivedMessage models.MessageCreate
	if err := json.Unmarshal(msg, &receivedMessage); err != nil {
		return err
	}

	// Vérifier si l'utilisateur est membre de l'association
	var membership models.Membership
	err := database.CurrentDatabase.
		Where("user_id = ? AND association_id = ? AND status = ?",
			user.ID, receivedMessage.AssociationID, enums.Accepted).
		First(&membership).Error
	if err != nil {
		return err
	}

	// Ajouter l'ID de l'expéditeur
	receivedMessage.SenderID = user.ID

	// Créer et diffuser le message
	createdMessage, err := s.messageService.CreateMessage(receivedMessage)
	if err != nil {
		return err
	}

	return s.BroadcastToAssociation(createdMessage.AssociationID, createdMessage)
}

// Envoyer un message WebSocket à une connexion donnée
func (s *WebSocketService) sendWebSocketMessage(data []byte, ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.TextMessage, data)
}

// Diffuser un message à toutes les connexions
func (s *WebSocketService) Broadcast(data []byte) error {
	for _, connection := range connections {
		if err := s.sendWebSocketMessage(data, connection.connection); err != nil {
			return err
		}
	}
	return nil
}

// Diffuser un message à une association spécifique
func (s *WebSocketService) BroadcastToAssociation(associationID string, message *models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, connection := range connections {
		// Charger les adhésions pour vérifier si l'utilisateur appartient à l'association
		var memberships []models.Membership
		if err := database.CurrentDatabase.Where("user_id = ? AND association_id = ?", connection.user.ID, associationID).Find(&memberships).Error; err != nil {
			return err
		}

		if len(memberships) > 0 { // Si des adhésions existent, l'utilisateur est membre de l'association
			if err := s.sendWebSocketMessage(data, connection.connection); err != nil {
				return err
			}
		}
	}

	return nil
}

// Gérer les messages d'une association
func (s *WebSocketService) FetchAssociationMessages(associationID string, ws *websocket.Conn) error {
	messages, err := s.messageService.GetMessagesByAssociation(associationID)
	if err != nil {
		return err
	}

	for _, message := range messages {
		data, err := json.Marshal(message)
		if err != nil {
			return err
		}
		if err := s.sendWebSocketMessage(data, ws); err != nil {
			return err
		}
	}

	return nil
}
