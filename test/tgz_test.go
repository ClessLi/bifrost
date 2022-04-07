package test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func tgz(dest string, list []string) (err error) {
	destDir := filepath.Dir(dest)
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tgzw := tar.NewWriter(gw)
	defer tgzw.Close()
	for _, f := range list {
		dirPath, dirErr := filepath.Rel(destDir, filepath.Dir(f))
		if dirErr != nil {
			return dirErr
		}
		fd, fErr := os.Open(f)
		if fErr != nil {
			return fErr
		}
		compErr := compress(fd, dirPath, tgzw)
		if compErr != nil {
			return compErr
		}
	}
	return
}

func compress(fd *os.File, prefix string, tgzw *tar.Writer) error {
	defer fd.Close()
	info, infoErr := fd.Stat()
	if infoErr != nil {
		return infoErr
	}
	header, hErr := tar.FileInfoHeader(info, "")
	if hErr != nil {
		return hErr
	}
	header.Name = filepath.Join(prefix, header.Name)
	twErr := tgzw.WriteHeader(header)
	if twErr != nil {
		return twErr
	}
	_, ioErr := io.Copy(tgzw, fd)
	fd.Close()
	if ioErr != nil {
		return ioErr
	}
	return nil
}

func TestTGZ(t *testing.T) {
	fmt.Println(os.Getwd())
	dest := `F:/Code_Path/src/bifrost/test/tgz_test/test.tgz`
	list := []string{
		`F:/Code_Path/src/bifrost/test/tgz_test/1/1.txt`,
		`F:/Code_Path/src/bifrost/test/tgz_test/2/2.txt`,
		`F:/Code_Path/src/bifrost/test/tgz_test/test.txt`,
	}
	err := tgz(dest, list)
	fmt.Println(os.Getwd())
	if err != nil {
		t.Log(err)
	}
}
