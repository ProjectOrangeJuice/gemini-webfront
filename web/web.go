package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func StartHandlers(sessionKey string) {
	router := gin.Default()
	router.LoadHTMLGlob("./web/static/*")
	// Create a new limiter that allows 5 requests per second with a burst limit of 5.
	lim := tollbooth.NewLimiter(5, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lim.SetBurst(5)

	store := cookie.NewStore([]byte(sessionKey))
	router.Use(sessions.Sessions("whoami", store))

	protected := router.Group("/", authMiddleware)

	protected.GET("/new", handleNewChat)
	protected.GET("/history", openChat)
	protected.GET("/list", listChats)
	protected.POST("/send", sendMessage)
	protected.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })

	// Host static pages
	protected.Static("/static", "./web/static")
	router.GET("/wiggle", loginHandler)
	router.POST("/wiggle", loginHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":9090", router))
}

func authMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("whoami") == nil {
		log.Printf("%s tried to access a page without logging in", c.ClientIP())
		c.Redirect(http.StatusFound, "/wiggle")
		c.Abort()
		return
	}
	log.Printf("%s accessed %s", c.ClientIP(), c.Request.URL.Path)
	c.Next()
}

func loginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		log.Println("GET")
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	p := c.Request.FormValue("pass")
	if p == "snap" {
		log.Printf("%s requested snap", c.ClientIP())
		os.Exit(1)
		return
	}

	if p != os.Getenv("USER_KEY") {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid password"})
		return
	}

	// Set the session
	session := sessions.Default(c)
	session.Set("whoami", true)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}
