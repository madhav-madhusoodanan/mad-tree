package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// FileItem represents a single file or folder for the JSON response
type FileItem struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

func BrowseFolder(c *gin.Context) {
	// The path the user wants to view (defaults to empty string for root)
	reqPath := c.Query("path")

	// SECURITY: Prevent directory traversal attacks (e.g., passing "../../etc")
	if strings.Contains(reqPath, "..") {
		c.JSON(400, gin.H{"error": "Invalid path requested"})
		return
	}

	// Join the requested path with the base public directory
	baseDir := "./public"
	fullPath := filepath.Join(baseDir, reqPath)

	// Ensure the public directory exists before reading
	os.MkdirAll(baseDir, 0777)

	// Read the directory contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(404, gin.H{"error": "Directory not found or is empty"})
		return
	}

	var items []FileItem
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Skip files we can't read
		}
		
		items = append(items, FileItem{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}

	c.JSON(200, gin.H{
		"currentPath": reqPath,
		"items":       items,
	})
}