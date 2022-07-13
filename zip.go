package gozip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ZipToFile(sourcePath, targetPath string) error {
	if sourcePath[len(sourcePath)-1] != '/' {
		sourcePath += "/"
	}

	var files []*FileToZip
	err := filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(sourcePath))

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		files = append(files, &FileToZip{
			Path:  relPath,
			Bytes: bytes,
		})
		return nil
	})
	if err != nil {
		return err
	}
	return ZipArray(files, targetPath)
}

type FileToZip struct {
	Path  string
	Bytes []byte
}

func ZipArray(files []*FileToZip, targetPath string) error {
	destinationFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	for _, file := range files {
		zipFile, err := myZip.Create(file.Path)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, bytes.NewReader(file.Bytes))
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func UnzipIoReadCloser(source io.ReadCloser, destination string) error {
	defer source.Close()

	allBytes, err := io.ReadAll(source)
	if err != nil {
		return err
	}

	bytesReader := bytes.NewReader(allBytes)

	reader, err := zip.NewReader(bytesReader, int64(len(allBytes)))
	if err != nil {
		return err
	}
	return unzipReader(reader, destination)
}

func ZipToIoReadCloser(files []*FileToZip) (io.ReadCloser, error) {
	buff := bytes.NewBuffer([]byte{})
	writer := zip.NewWriter(buff)
	for _, file := range files {
		zipFile, err := writer.Create(file.Path)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(zipFile, bytes.NewReader(file.Bytes))
		if err != nil {
			return nil, err
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(buff), nil
}
