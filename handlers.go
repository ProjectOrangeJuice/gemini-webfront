package main

import (
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

	log.Fatal(http.ListenAndServe(":9090", router))
}

const chatDirectory = "chats"

func handleNewChat(g *gin.Context) {
	// Generate a new chat file, return the name

	// Random 6 char string
	token := uuid.New()

	// Create the file
	file, err := os.Create(chatDirectory + "/" + token.String() + ".json")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

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

	list := make([]string, len(files))
	for i, f := range files {
		list[i] = strings.Replace(f.Name(), ".json", "", 1)
	}

	// return the list of files
	g.JSON(200, gin.H{
		"chats": list,
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
		"message": resp,
	})
}
