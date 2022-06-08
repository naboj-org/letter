package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"naboj.org/letter/runner"
	"net/http"
	"os"
	"time"
)

func GenerateAsynchronously(c *gin.Context) {
	r := validateRequest(c, true)
	if !r.Ok {
		if r.File != nil {
			r.File.Close()
		}
		return
	}

	defer r.File.Close()

	go workAsync(r.File, r.Entrypoint, r.Callback)
	c.Status(http.StatusCreated)
}

func workAsync(r io.Reader, entrypoint, returnUrl string) {
	log.Println(returnUrl)
	o, err := runner.ProcessArchive(r, entrypoint)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	fw, _ := w.CreateFormField("tex_output")
	if err != nil {
		log.Println("create form field 'tex_output' failed:", err)
		return
	}
	fw.Write([]byte(o.Output))

	if err != nil {
		fw, err := w.CreateFormField("error")
		if err != nil {
			log.Println("create form field 'error' failed:", err)
			return
		}
		fw.Write([]byte(err.Error()))
	} else {
		if o.File != "" {
			defer os.Remove(o.File)

			fw, err := w.CreateFormFile("file", "output.pdf")
			if err != nil {
				log.Println("create form file failed:", err)
				return
			}

			file, err := os.Open(o.File)
			if err != nil {
				log.Println("opening output file failed:", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(fw, file)
			if err != nil {
				log.Println("copying file to request failed:", err)
				return
			}
		} else {
			fw, err := w.CreateFormField("error")
			if err != nil {
				log.Println("create form field 'error' failed:", err)
				return
			}
			fw.Write([]byte("TeX did not generate any file"))
		}
	}

	w.Close()
	req, err := http.NewRequest("POST", returnUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Println("creating request failed:", err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = client.Do(req)
	if err != nil {
		log.Println("request error:", err)
	}
}
