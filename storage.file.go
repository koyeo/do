package do

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	Path     string
	Name     string
	Size     int64
	MimeType string
	Md5      string
	Url      string
}

func ChooseFile(path string) (file *File, err error) {

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	path = filepath.Join(dir, path)
	info, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			err = fmt.Errorf("file '%s' not exist", path)
			return
		}
		err = fmt.Errorf("check file '%s' error: %s", path, err.Error())
		return
	}

	if info.IsDir() {
		err = fmt.Errorf("'%s' is a dir, expect a file", path)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	file = new(File)
	file.Path = path
	file.Name = info.Name()
	file.Size = info.Size()
	file.MimeType, err = fileMimeType(f)
	if err != nil {
		return
	}
	file.Md5 = fileMd5(f)

	return
}

func fileMimeType(out *os.File) (string, error) {

	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func fileMd5(f *os.File) string {

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
