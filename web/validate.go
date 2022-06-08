package web

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strings"
)

type requestData struct {
	Ok         bool
	Entrypoint string
	Callback   string
	File       multipart.File
}

func validateRequest(c *gin.Context, hasCallback bool) requestData {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return requestData{Ok: false}
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return requestData{Ok: false}
	}

	entrypoint, present := c.GetPostForm("entrypoint")
	if !present || strings.TrimSpace(entrypoint) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entrypoint not specified"})
		return requestData{Ok: false, File: file}
	}

	callback := ""
	if hasCallback {
		callback, present = c.GetPostForm("callback")
		if !present || strings.TrimSpace(callback) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "callback not specified"})
			return requestData{Ok: false, File: file}
		}
	}

	return requestData{
		Ok:         true,
		Entrypoint: entrypoint,
		Callback:   callback,
		File:       file,
	}
}
