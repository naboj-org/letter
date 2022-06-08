package untar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
)

func UnTar(r io.Reader, root string) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeDir {
			err := os.Mkdir(path.Join(root, header.Name), os.ModePerm)
			if err != nil {
				return err
			}
		} else if header.Typeflag == tar.TypeReg {
			file, err := os.Create(path.Join(root, header.Name))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
			file.Close()
		} else {
			return fmt.Errorf("unknown tar type: %v", header.Typeflag)
		}
	}

	return nil
}
