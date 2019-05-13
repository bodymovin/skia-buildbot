package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	assert "github.com/stretchr/testify/require"
	"go.skia.org/infra/go/testutils"
	"go.skia.org/infra/go/testutils/unittest"
)

func createFile(dir, prefix, content string, t *testing.T) string {
	f, err := ioutil.TempFile(dir, prefix)
	assert.NoError(t, err)
	_, err = f.WriteString(content)
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
	return f.Name()
}

func assertFileExists(dir, path, content string, t *testing.T) {
	c, err := ioutil.ReadFile(filepath.Join(dir, filepath.Base(path)))
	assert.NoError(t, err)
	assert.Equal(t, content, string(c))
}

func TestZipE2E(t *testing.T) {
	unittest.MediumTest(t)

	// Create a directory in temp.
	targetDir, err := ioutil.TempDir("", "zip_test")
	assert.NoError(t, err)
	defer testutils.RemoveAll(t, targetDir)

	// Populate the target dir.
	// Create two files in target dir.
	f1 := createFile(targetDir, "temp1", "testing1", t)
	f2 := createFile(targetDir, "temp2", "testing2", t)
	// Create subdir.
	subDir, err := ioutil.TempDir(targetDir, "zip_test")
	assert.NoError(t, err)
	// Create one file in subdir.
	f3 := createFile(subDir, "temp3", "testing3", t)
	// Create an empty subdir.
	emptySubDir, err := ioutil.TempDir(targetDir, "zip_test")
	assert.NoError(t, err)

	// Zip the directory.
	outputDir, err := ioutil.TempDir("", "zip_location")
	defer testutils.RemoveAll(t, outputDir)
	zipPath := filepath.Join(outputDir, "test.zip")
	err = ZipIt(zipPath, targetDir)
	assert.NoError(t, err)
	// Assert that zip was created
	_, err = os.Stat(zipPath)
	assert.NoError(t, err)

	// Test UnZipping.
	err = UnZip(outputDir, zipPath)
	assert.NoError(t, err)
	// Assert that the 3 zipped files are in the right locations.
	assertFileExists(outputDir, f1, "testing1", t)
	assertFileExists(outputDir, f2, "testing2", t)
	assertFileExists(filepath.Join(outputDir, filepath.Base(subDir)), f3, "testing3", t)
	// Assert that the empty subdir was created.
	_, err = os.Stat(filepath.Join(outputDir, filepath.Base(emptySubDir)))
	assert.NoError(t, err)
}
