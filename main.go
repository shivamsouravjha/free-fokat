package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

func startTicker() *time.Ticker {
	ticker := time.NewTicker(48 * time.Second)

	go func() {
		for t := range ticker.C {

			cmd := exec.Command("curl", "https://stock-backend-3b3b.onrender.com/api/keepServerRunning")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return
			}
			fmt.Println(string(output), "output", t)
			cmd = exec.Command("curl", "https://emogpt.onrender.com/api/keepServerRunning")
			output, err = cmd.CombinedOutput()
			if err != nil {
				return
			}
			fmt.Println(string(output), "output", t)
		}
	}()

	return ticker
}

func IsRunning(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Server is running"})
}

func Routes(r *gin.Engine) {

	v1 := r.Group("/api")

	{
		v1.GET("/keepServerRunning", IsRunning)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic and stack trace
				// Respond with a 500 Internal Server Error
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error. Please try again later.",
				})
				ctx.Abort()
			}
		}()
		// Continue to the next handler
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, trell-auth-token, trell-app-version-int, creator-space-auth-token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.New()
	router.Use(RecoveryMiddleware())

	router.Use(CORSMiddleware())

	startTicker()
	Routes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4002"
	}

	// Create a server instance using gin engine as handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

}
