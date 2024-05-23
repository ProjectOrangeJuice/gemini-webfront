package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func startHandlers() {
	router := gin.Default()
	// Create a new limiter that allows 5 requests per second with a burst limit of 5.
	lim := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lim.SetBurst(5)

	router.GET("/new", handleNewChat)
	router.GET("/history", openChat)
	router.GET("/list", listChats)
	router.POST("/send", sendMessage)

	// Host static pages
	router.Static("/static", "./static")

	// Start the server

	log.Fatal(http.ListenAndServe(":9090", router))
}

const chatDirectory = "chats"

func handleNewChat(g *gin.Context) {
	// Generate a new chat file, return the name

	// Random 6 char string
	token := uuid.New()

	chatHistory := chatFile{
		Title: "untitled",
		When:  time.Now(),
	}

	bytes, err := json.Marshal(chatHistory)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s.json", chatDirectory, token), bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// return the token to the user
	g.JSON(200, gin.H{
		"token": token.String(),
	})
}

func listChats(g *gin.Context) {
	// Get a list of files in a directory
	files, err := os.ReadDir(chatDirectory)
	if err != nil {
		log.Fatal(err)
	}

	c := make([]chats, len(files))
	for i, f := range files {
		chatFile := readChat(strings.Replace(f.Name(), ".json", "", 1))
		c[i] = chats{
			Title: chatFile.Title,
			When:  chatFile.When,
			Token: strings.Replace(f.Name(), ".json", "", 1),
		}
	}

	// return the list of files
	g.JSON(200, gin.H{
		"Chats": c,
	})
}

func openChat(g *gin.Context) {
	// Get the chat file
	token := g.Query("token")
	g.File(chatDirectory + "/" + token + ".json")
}

func sendMessage(g *gin.Context) {
	token := g.Query("token")
	message := g.PostForm("message")
	resp, _ := sendChat(token, message)
	g.JSON(200, gin.H{
		"Message": resp,
	})
}
