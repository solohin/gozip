package gozip_test

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/solohin/gozip"
	"github.com/stretchr/testify/require"
)

func TestZipFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.MkdirAll(dir+"/a/b/c", 0755)
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/a/b/c/d.txt", []byte("abcd"), 0644)
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/root.txt", []byte("root"), 0644)
	require.NoError(t, err)

	zipFilePath, err := ioutil.TempFile("", "*.zip")
	defer os.Remove(zipFilePath.Name())
	require.NoError(t, err)

	err = gozip.ZipToFile(dir, zipFilePath.Name())
	require.NoError(t, err)

	err = gozip.UnzipFromFile(zipFilePath.Name(), dir+"/unzipped")
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile(dir + "/unzipped/root.txt")
	require.NoError(t, err)
	require.Equal(t, "root", string(bytes))

	bytes, err = ioutil.ReadFile(dir + "/unzipped/a/b/c/d.txt")
	require.NoError(t, err)
	require.Equal(t, "abcd", string(bytes))
}

func TestZipToIoReadCloser(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	var files []*gozip.FileToZip

	files = append(files, &gozip.FileToZip{
		Path:  "/a/b/c/d.txt",
		Bytes: []byte("abcd"),
	})

	files = append(files, &gozip.FileToZip{
		Path:  "/root.txt",
		Bytes: []byte("root"),
	})

	var result io.ReadCloser
	result, err = gozip.ZipToIoReadCloser(files)
	require.NoError(t, err)

	zipFilePath, err := ioutil.TempFile("", "*.zip")
	defer os.Remove(zipFilePath.Name())
	require.NoError(t, err)

	outFile, err := os.Create(zipFilePath.Name())
	require.NoError(t, err)
	defer outFile.Close()

	_, err = io.Copy(outFile, result)
	require.NoError(t, err)

	err = gozip.UnzipFromFile(zipFilePath.Name(), dir+"/unzipped")
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile(dir + "/unzipped/root.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "root", string(bytes))

	bytes, err = ioutil.ReadFile(dir + "/unzipped/a/b/c/d.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "abcd", string(bytes))
}

func TestUnzipIoReadCloser(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.MkdirAll(dir+"/a/b/c", 0755)
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/a/b/c/d.txt", []byte("abcd"), 0644)
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/root.txt", []byte("root"), 0644)
	require.NoError(t, err)

	zipFilePath, err := ioutil.TempFile("", "*.zip")
	defer os.Remove(zipFilePath.Name())
	require.NoError(t, err)

	err = gozip.ZipToFile(dir, zipFilePath.Name())
	require.NoError(t, err)

	reader, err := os.Open(zipFilePath.Name())
	require.NoError(t, err)

	err = gozip.UnzipIoReadCloser(reader, dir+"/unzipped")
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile(dir + "/unzipped/root.txt")
	require.NoError(t, err)
	require.Equal(t, "root", string(bytes))

	bytes, err = ioutil.ReadFile(dir + "/unzipped/a/b/c/d.txt")
	require.NoError(t, err)
	require.Equal(t, "abcd", string(bytes))
}

func TestIoReadCloserCycle(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	var files []*gozip.FileToZip

	files = append(files, &gozip.FileToZip{
		Path:  "/a/b/c/d.txt",
		Bytes: []byte("abcd"),
	})

	files = append(files, &gozip.FileToZip{
		Path:  "/root.txt",
		Bytes: []byte("root"),
	})

	var result io.ReadCloser
	result, err = gozip.ZipToIoReadCloser(files)
	require.NoError(t, err)

	err = gozip.UnzipIoReadCloser(result, dir+"/unzipped")
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile(dir + "/unzipped/root.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "root", string(bytes))

	bytes, err = ioutil.ReadFile(dir + "/unzipped/a/b/c/d.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "abcd", string(bytes))
}

func debugTree(t *testing.T, dir string) {
	cmd := exec.Command("tree", dir)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log(string(out))
}
