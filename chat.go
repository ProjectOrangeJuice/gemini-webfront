package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	currentToken  string
	currentClient *genai.Client
	currentSender func(msg string) *genai.GenerateContentResponse
)

func sendChat(token string, message string) (string, error) {
	log.Printf("message for %s", token)
	if token != currentToken {
		log.Println("Starting new chat")
		startNewChat(token)
	}

	resp := currentSender(message)
	msg := readResponse(resp)

	appendToHistory(token, message, msg)
	return msg, nil
}

func appendToHistory(token, message, response string) {
	// Read the current file
	// Read the chats from the file
	var chatHistory chatFile
	// read the file
	file, err := os.Open(chatDirectory + "/" + token + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &chatHistory)
	if err != nil {
		log.Fatal(err)
	}

	// append messages
	chatHistory.Messages = append(chatHistory.Messages, chat{
		Content: message,
		Who:     "user",
	})

	chatHistory.Messages = append(chatHistory.Messages, chat{
		Content: response,
		Who:     "model",
	})

	// overwrite the file with the new json
	bytes, err := json.Marshal(chatHistory)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(chatDirectory+"/"+token+".json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func startNewChat(token string) {
	// end previous
	if currentClient != nil {
		currentClient.Close()
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	currentSender = func(msg string) *genai.GenerateContentResponse {
		fmt.Printf("== Me: %s\n== Model:\n", msg)
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	currentToken = token
	currentClient = client

	// Populate the history

	// Read the chats from the file
	chatFile := readChat(token)
	for _, chat := range chatFile.Messages {
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(chat.Content),
			},
			Role: chat.Who,
		})
	}
}

func readChat(token string) chatFile {
	var chat chatFile
	// read the file
	file, err := os.Open(chatDirectory + "/" + token + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &chat)
	if err != nil {
		log.Fatal(err)
	}
	return chat
}
