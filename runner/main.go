package runner

import (
	"bytes"
	"io"
	"io/ioutil"
	"naboj.org/letter/untar"
	"os"
	"os/exec"
	"path"
)

func RunLatex(directory, filename string) (string, error) {
	cmd := exec.Command("/usr/bin/latexmk", "-lualatex", "-norc", "-jobname=lttr_output", filename)
	cmd.Dir = directory
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	return b.String(), err
}

type ProcessResult struct {
	Output string
	File   string
}

func ProcessArchive(r io.Reader, filename string) (ProcessResult, error) {
	dir, err := ioutil.TempDir("", "lttr")
	if err != nil {
		return ProcessResult{}, err
	}
	defer os.RemoveAll(dir)

	err = untar.UnTar(r, dir)
	if err != nil {
		return ProcessResult{}, err
	}

	out, err := RunLatex(dir, filename)
	if err != nil {
		return ProcessResult{Output: out}, err
	}

	old, err := os.Open(path.Join(dir, "lttr_output.pdf"))
	if err != nil {
		return ProcessResult{}, err
	}
	defer old.Close()

	newFile, err := ioutil.TempFile("", "lttrout")
	if err != nil {
		return ProcessResult{}, err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, old)
	if err != nil {
		return ProcessResult{}, err
	}

	return ProcessResult{File: newFile.Name(), Output: out}, nil
}
