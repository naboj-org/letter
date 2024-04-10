package pkg

import (
	"code.cloudfoundry.org/archiver/extractor"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func unzip(zipfile io.Reader, directory string) error {
	tempZipName := path.Join(directory, "_source.zip")
	zip, err := os.Create(tempZipName)
	if err != nil {
		return err
	}
	defer os.Remove(tempZipName)

	_, err = io.Copy(zip, zipfile)
	if err != nil {
		return err
	}

	err = extractor.NewZip().Extract(zip.Name(), directory)
	if err != nil {
		return err
	}
	return nil
}

func ProcessJob(zipfile io.Reader, entrypoint string) (string, string, error) {
	dir, err := os.MkdirTemp("", "lttr")
	if err != nil {
		return "", "", fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(dir)

	err = unzip(zipfile, dir)
	if err != nil {
		return "", "", fmt.Errorf("extracting zip: %w", err)
	}

	out, err := RunTectonic(dir, entrypoint)
	if err != nil {
		return out, "", fmt.Errorf("tectonic: %w", err)
	}

	outputName := strings.TrimSuffix(entrypoint, ".tex") + ".pdf"
	old, err := os.Open(path.Join(dir, outputName))
	if err != nil {
		return out, "", fmt.Errorf("opening output: %w", err)
	}
	defer old.Close()

	newFile, err := os.CreateTemp("", "lttrout")
	if err != nil {
		return out, "", fmt.Errorf("creating final file: %w", err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, old)
	if err != nil {
		return out, "", fmt.Errorf("copying output: %w", err)
	}

	return out, newFile.Name(), nil
}
