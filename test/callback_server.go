package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		log.Println("err:", c.PostForm("error"))
		log.Println("log:", c.PostForm("tectonic_output"))

		fileHeader, err := c.FormFile("file")
		if err != nil {
			log.Println("ERROR", err)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			log.Println("ERROR", err)
			return
		}

		store, err := os.Create("output.pdf")
		if err != nil {
			log.Println("ERROR", err)
			return
		}

		_, err = io.Copy(store, file)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	})
	r.Run(":8081")
}
