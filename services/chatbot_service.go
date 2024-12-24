package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatService représente le service qui interagit avec l'API OpenAI
type ChatService struct {
}

// NewChatService crée une nouvelle instance de ChatService avec la clé API
func NewChatService() *ChatService {
	return &ChatService{}
}

// ChatRequest représente la requête envoyée au chatbot
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse représente la réponse du service
type ChatResponse struct {
	Response string `json:"response"`
}

// OpenAIResponse représente la structure de la réponse de l'API OpenAI
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GetChatGPTResponse envoie un message à l'API OpenAI et retourne la réponse
func (cs *ChatService) GetChatGPTResponse(message string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo", // ou "gpt-4" selon l'accès
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %v", err)
	}

	// Crée la requête HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cs.APIKey))

	// Exécute la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to OpenAI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status)
	}

	// Lit et décode la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result OpenAIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Vérifie et renvoie le contenu de la réponse
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("empty response from ChatGPT")
}
