package main

import (
	"log"
	"os"

	"github.com/ProjectOrangeJuice/gemini-webfront/gemini"
	"github.com/ProjectOrangeJuice/gemini-webfront/web"
)

// Write a chatbot that lets me have a converstation with gemini
func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing api key")
	}

	gemini.GeminiKey = apiKey
	web.StartHandlers()
}
