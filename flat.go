package gozip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzipFileFlat(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	fileName := filepath.Base(f.Name)

	filePath := filepath.Join(destination, fileName)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		return nil
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func unzipReaderFlat(reader *zip.Reader, destination string) error {
	// 2. Get the absolute destination path
	destination, err := filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFileFlat(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnzipFromBytesFlat(zipBytes []byte, target string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return err
	}

	return unzipReaderFlat(zipReader, target)
}

func UnzipFromFileFlat(sourcePath, targetPath string) error {
	reader, err := zip.OpenReader(sourcePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	return unzipReaderFlat(&reader.Reader, targetPath)
}
