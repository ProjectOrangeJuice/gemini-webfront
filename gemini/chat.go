package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const ChatDirectory = "chats"

func startNewChat(token string) error {
	// end previous
	if activeChat.client != nil {
		activeChat.client.Close()
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiKey))
	if err != nil {
		return fmt.Errorf("could not start a new chat, %w", err)
	}

	model := client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	currentSender := func(msg string) *genai.GenerateContentResponse {
		log.Printf("== Me: %s\n", msg)
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	activeChat.token = token
	activeChat.client = client
	activeChat.sender = currentSender
	// Populate the history

	// Read the chats from the file
	chatFile, err := ReadChat(token)
	if err != nil {
		activeChat.newChat = true
		log.Printf("could not load in history, %v")
		return nil
	}
	if len(chatFile.Messages) == 0 {
		activeChat.newChat = true
	}

	for _, chat := range chatFile.Messages {
		activeChat.newChat = false
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(chat.Content),
			},
			Role: chat.Who,
		})
	}
	return nil
}

func ReadChat(token string) (chatFile, error) {
	var chat chatFile
	// read the file
	file, err := os.Open(ChatDirectory + "/" + token + ".json")
	if err != nil {
		return chat, fmt.Errorf("could not open chatfile, %w", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return chat, fmt.Errorf("could not read chatfile, %w", err)
	}
	err = json.Unmarshal(content, &chat)
	if err != nil {
		return chat, fmt.Errorf("could not unmarshal chatfile, %w", err)
	}
	return chat, nil
}

func ListChats() []chats {
	// Get a list of files in a directory
	files, err := os.ReadDir(ChatDirectory)
	if err != nil {
		log.Printf("Can't read directory, %s", err)
		return nil
	}

	c := make([]chats, len(files))
	for i, f := range files {
		chatFile, err := ReadChat(strings.Replace(f.Name(), ".json", "", 1))
		if err != nil {
			log.Printf("read chat for listing, %s", err)
			continue
		}
		c[i] = chats{
			Title: chatFile.Title,
			When:  chatFile.When,
			Token: strings.Replace(f.Name(), ".json", "", 1),
		}
	}
	// sort the chats by date
	sort.Slice(c, func(i, j int) bool {
		return c[i].When.Before(c[j].When)
	})

	return c
}
