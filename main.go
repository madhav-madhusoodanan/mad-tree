package main

import (
  "log"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/madhav-madhusoodanan/mad-tree/internal"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
   	r.MaxMultipartMemory = 8 << 20 // 8 MiB
    r.Static("/assets", "./public")
  
    // 1. Tell Gin where to find your HTML templates
    r.LoadHTMLGlob("templates/*")

    // Define a simple GET endpoint
    r.GET("/ping", func(c *gin.Context) {
    	// Return JSON response
	    c.JSON(http.StatusOK, gin.H{
	      "message": "pong",
	    })
    })
  
    // 2. Serve the HTML file on the root route
	r.GET("/", func(c *gin.Context) {
	    c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/browse", func(c *gin.Context) { 
		c.HTML(http.StatusOK, "browse.html", nil) 
	})
	r.GET("/api/browse", internal.BrowseFolder)
  
  	r.POST("/upload/:name", internal.UploadFile)

   	// Start server on port 8080 (default)
    // Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
    if err := r.Run("0.0.0.0:8080"); err != nil {
    	log.Fatalf("failed to run server: %v", err)
    }
}