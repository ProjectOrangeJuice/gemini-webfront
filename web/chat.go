package web

import (
	"github.com/ProjectOrangeJuice/gemini-webfront/gemini"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleNewChat(g *gin.Context) {
	// Generate a new chat file, return the name

	// Random 6 char string
	token := uuid.New()
	// return the token to the user
	g.JSON(200, gin.H{
		"Token": token.String(),
	})
}

func listChats(g *gin.Context) {
	chats := gemini.ListChats()
	// return the list of files
	g.JSON(200, gin.H{
		"Chats": chats,
	})
}

func openChat(g *gin.Context) {
	// Get the chat file
	token := g.Query("token")
	g.File(gemini.ChatDirectory + "/" + token + ".json")
}

func sendMessage(g *gin.Context) {
	token := g.Query("token")
	message := g.PostForm("message")
	resp, err := gemini.SendChat(token, message)
	if err != nil {
		g.AbortWithError(500, err)
		return
	}
	g.JSON(200, gin.H{
		"Message": resp,
	})
}
