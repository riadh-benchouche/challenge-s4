package services

import (
	"backend/database"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (cs *ChatService) GetChatGPTResponse(message string, dbData []models.AssociationSummary) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Formatez les données de la base à intégrer
	dbInfo := "Voici les associations disponibles en lien avec votre requête:\n"
	for _, assoc := range dbData {
		dbInfo += fmt.Sprintf("- %s : %s\n", assoc.Name, assoc.Description)
	}

	// Prépare la requête
	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Tu es un assistant pour les étudiants de l'ESGI, qui répond selon les données d'association fournies.",
			},
			{
				"role":    "assistant",
				"content": dbInfo,
			},
			{
				"role":    "user",
				"content": message,
			},
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %v", err)
	}

	// HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key is missing. Please set the OPENAI_API_KEY environment variable")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result OpenAIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty response from ChatGPT")
}

func (as *AssociationService) SearchAssociations() ([]models.AssociationSummary, error) {
	var associations []models.AssociationSummary
	if err := database.CurrentDatabase.
		Table("associations").
		Select("name, description").
		Find(&associations).Error; err != nil {
		return nil, err
	}
	return associations, nil
}
