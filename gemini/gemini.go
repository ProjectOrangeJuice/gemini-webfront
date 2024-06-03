package gemini

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// Only building this to run one chat at a time

var (
	activeChat liveChat
	GeminiKey  string
)

type liveChat struct {
	client  *genai.Client
	token   string
	newChat bool
	sender  func(msg string) *genai.GenerateContentResponse
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
	AddToHistory(token, chat{
		Who:     "user",
		Content: message,
	})
	resp := activeChat.sender(message)
	respMsg := readResponse(resp)
	AddToHistory(token, chat{
		Who:     "model",
		Content: respMsg,
	})

	if activeChat.newChat {
		activeChat.newChat = false
		go generateTitleBackground(token)
	}

	return respMsg, nil

}

func generateTitleBackground(token string) {
	resp := activeChat.sender("For the next message and just the next message you're talking to a program. It's looking for a title to give this chat. Reply in no more than 10 words what this chat should be called.")
	respMsg := readResponse(resp)
	respMsg = strings.ReplaceAll(respMsg, "**", "")
	ChangeTitle(token, respMsg)
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
	log.Printf("Model: %s", response)
	return response // Return the concatenated string
}
