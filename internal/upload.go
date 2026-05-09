package internal

import (
  "io"
  "os"
  "log/slog"
  "github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	/*
	 *
		* 1. Fetch the data stream
		* 2. Iterate through all the chunks
		* 3. Write to a file in local
		* 
	 */
	logger := slog.Default()
	reader, err := c.Request.MultipartReader()
	if err != nil {
		return
	}
	
	folderName := c.Param("name")
	completeFolderName := "./public/" + folderName
	fileNames := []string{}

	err = os.MkdirAll(completeFolderName, 0777)
	if err != nil {
		logger.Error("Error in creating folder (" + completeFolderName + "): " + err.Error())
		return
	}

	for {
		part, err := reader.NextPart()

		if err == io.EOF {
			logger.Debug("EOF detected!")
			break
		}
		
		if part.FileName() == "" {
			continue
		}

		fileName := completeFolderName + "/" + part.FileName()
		out, err := os.Create(fileName)
		
		if err != nil {
			logger.Error("Error in creating file: " + fileName + ": " + err.Error())
			continue
		}

		buffer := make([]byte, 2*1024*1024) // 1 MB buffer
		_, err = io.CopyBuffer(out, part, buffer)
		out.Close()

		if err != nil {
			logger.Error("Error in uploading " + fileName + ": " + err.Error())
			continue
		}

		fileNames = append(fileNames, fileName)
	}
	
	c.JSON(200, gin.H{"status": "Files received", "names": fileNames})
}