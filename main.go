package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"naboj.org/letter/web"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	token, present := os.LookupEnv("AUTH_TOKEN")
	if !present {
		log.Println("AUTH_TOKEN is not configured, authentication is disabled.")
	} else {
		r.Use(web.AuthHandler(token))
	}

	r.POST("/sync", web.GenerateSynchronously)
	r.POST("/async", web.GenerateAsynchronously)
	r.Run()
}
