package testutils

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func MakeRandomTestFile(t *testing.T, size int64, copyTo io.Writer) string {
	file, err := ioutil.TempFile("", "destructive_cat_test.*.data")
	AssertNoErr(t, err, "cannot create test file")
	name := file.Name()
	t.Cleanup(func() {
		_ = os.Remove(name)
	})
	writer := io.Writer(file)
	if copyTo != nil {
		writer = io.MultiWriter(file, copyTo)
	}
	_, err = io.CopyN(writer, rand.Reader, size)
	AssertNoErr(t, err, "cannot write test data")
	err = file.Close()
	AssertNoErr(t, err, "cannot close test file")
	return name
}
