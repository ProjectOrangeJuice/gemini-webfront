package web

import (
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

func StartHandlers() {
	router := gin.Default()
	// Create a new limiter that allows 5 requests per second with a burst limit of 5.
	lim := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lim.SetBurst(5)

	router.GET("/new", handleNewChat)
	router.GET("/history", openChat)
	router.GET("/list", listChats)
	router.POST("/send", sendMessage)

	// Host static pages
	router.Static("/static", "./web/static")

	// Start the server

	log.Fatal(http.ListenAndServe(":9090", router))
}
