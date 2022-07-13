package gozip_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/solohin/gozip"
	"github.com/stretchr/testify/require"
)

func TestFlat(t *testing.T) {
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

	bytes, err := ioutil.ReadAll(result)
	require.NoError(t, err)

	err = os.MkdirAll(dir+"/unzipped", 0755)
	require.NoError(t, err)

	err = gozip.UnzipFromBytesFlat(bytes, dir+"/unzipped")
	require.NoError(t, err)

	bytes, err = ioutil.ReadFile(dir + "/unzipped/root.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "root", string(bytes))

	bytes, err = ioutil.ReadFile(dir + "/unzipped/d.txt")
	if err != nil {
		debugTree(t, dir)
	}
	require.NoError(t, err)
	require.Equal(t, "abcd", string(bytes))
}
