package pkg

import (
	"bytes"
	"os/exec"
)

func RunTectonic(workDirectory string, entrypoint string) (string, error) {
	cmd := exec.Command("/usr/bin/tectonic", "-X", "compile", entrypoint)
	cmd.Dir = workDirectory
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	err := cmd.Run()
	return b.String(), err
}
