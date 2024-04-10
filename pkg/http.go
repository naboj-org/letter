package pkg

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func checkToken(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if token != authToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied."})
			c.Abort()
			return
		}
	}
}

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	token, present := os.LookupEnv("AUTH_TOKEN")
	if !present {
		log.Println("AUTH_TOKEN is not configured, authentication is disabled.")
	} else {
		r.Use(checkToken(token))
	}

	r.POST("/sync", handleSynchronous)
	r.POST("/async", handleAsynchronous)
	return r
}

func handleSynchronous(c *gin.Context) {
	r := validateRequest(c, false)
	if r.File != nil {
		defer r.File.Close()
	}
	if !r.Ok {
		return
	}

	output, file, err := ProcessJob(r.File, r.Entrypoint)
	if file != "" {
		defer os.Remove(file)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":           err.Error(),
			"tectonic_output": output,
		})
		return
	}

	c.File(file)
}

func handleAsynchronous(c *gin.Context) {
	r := validateRequest(c, true)
	if !r.Ok {
		if r.File != nil {
			r.File.Close()
		}
		return
	}

	go workAsync(r.File, r.Entrypoint, r.Callback)
	c.Status(http.StatusCreated)
}
