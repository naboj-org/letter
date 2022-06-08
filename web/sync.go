package web

import (
	"github.com/gin-gonic/gin"
	"naboj.org/letter/runner"
	"net/http"
	"os"
)

func GenerateSynchronously(c *gin.Context) {
	r := validateRequest(c, false)
	if !r.Ok {
		if r.File != nil {
			r.File.Close()
		}
		return
	}

	defer r.File.Close()

	o, err := runner.ProcessArchive(r.File, r.Entrypoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      err.Error(),
			"tex_output": o.Output,
		})
		return
	}

	if o.File == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "TeX did not generate any file",
			"tex_output": o.Output,
		})
		return
	}

	defer os.Remove(o.File)
	c.File(o.File)
}
