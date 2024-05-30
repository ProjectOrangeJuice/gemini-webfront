package gemini

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type chatFile struct {
	When     time.Time
	Title    string
	Messages []chat
}

type chat struct {
	Who     string
	Content string
}

type chats struct {
	Token string
	Title string
	When  time.Time
}

func AddToHistory(token string, msg chat) {
	chatHistory, err := ReadChat(token)
	if err != nil {
		log.Printf("can't update the history, %v", err)
	}

	chatHistory.Messages = append(chatHistory.Messages, msg)
	bytes, err := json.Marshal(chatHistory)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(ChatDirectory+"/"+token+".json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
