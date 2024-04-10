package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func workAsync(zipfile io.ReadCloser, entrypoint string, callbackUrl string) {
	defer zipfile.Close()

	output, file, err := ProcessJob(zipfile, entrypoint)
	if file != "" {
		defer os.Remove(file)
	}

	err = callback(callbackUrl, output, err, file)
	if err != nil {
		log.Printf("Could not send callback to %s: %v\n", callbackUrl, err)
	}
}

func callback(url string, tectonicOutput string, tectonicErr error, outputFile string) error {
	if outputFile != "" {
		defer os.Remove(outputFile)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	defer w.Close()

	fw, err := w.CreateFormField("tectonic_output")
	if err != nil {
		return fmt.Errorf("create tectonic_output field: %w", err)
	}
	_, err = fw.Write([]byte(tectonicOutput))
	if err != nil {
		return fmt.Errorf("populate tectonic_output field: %w", err)
	}

	if tectonicErr != nil {
		fw, err := w.CreateFormField("error")
		if err != nil {
			return fmt.Errorf("create error field: %w", err)
		}
		_, err = fw.Write([]byte(tectonicErr.Error()))
		if err != nil {
			return fmt.Errorf("populate error field: %w", err)
		}
	}

	if outputFile != "" {
		fw, err := w.CreateFormFile("file", "output.pdf")
		if err != nil {
			return fmt.Errorf("create file field: %w", err)
		}

		file, err := os.Open(outputFile)
		if err != nil {
			return fmt.Errorf("open output file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(fw, file)
		if err != nil {
			return fmt.Errorf("copy output to request: %w", err)
		}
	}

	w.Close()

	req, err := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	return nil
}
