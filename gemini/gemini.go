package gemini

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// Only building this to run one chat at a time

var (
	activeChat liveChat
	GeminiKey  string
)

type liveChat struct {
	client *genai.Client
	token  string
	sender func(msg string) *genai.GenerateContentResponse
}

func SendChat(token, message string) (string, error) {
	if message == "" {
		return "", nil
	}

	// check if there's an active chat
	if activeChat.client != nil || activeChat.token != token {
		err := startNewChat(token)
		if err != nil {
			return "", fmt.Errorf("can't send message, %s", err)
		}
	}

	resp := activeChat.sender(message)
	return readResponse(resp), nil

}

func readResponse(resp *genai.GenerateContentResponse) string { // Return type is now string
	response := "" // Initialize response variable
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response += fmt.Sprintf("%s", part)
			}
		}
	}
	return response // Return the concatenated string
}
