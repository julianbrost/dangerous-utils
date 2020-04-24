package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestEmptyFiles(t *testing.T) {
	output := bytes.NewBuffer(nil)
	err := DestructiveCat([]string{}, NopCloser(output), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != 0 {
		t.Errorf("DestructiveCat with an empty file list wrote %d bytes", output.Len())
	}
}

func TestSimple(t *testing.T) {
	output := bytes.NewBuffer(nil)
	copy := bytes.NewBuffer(nil)
	file := makeRandomTestFile(t, DefaultBlockSize, copy)
	err := DestructiveCat([]string{file}, NopCloser(output), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != copy.Len() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.Len(), copy.Len())
	}
	if bytes.Compare(output.Bytes(), copy.Bytes()) != 0 {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestInputDeleted(t *testing.T) {
	size := int64(1024)
	file := makeRandomTestFile(t, size, nil)
	err := DestructiveCat([]string{file}, NopCloser(ioutil.Discard), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		t.Error("input file was not deleted")
	}
}

func TestMultiBlock(t *testing.T) {
	output := NewCountingHashWriter()
	copy := NewCountingHashWriter()
	file := makeRandomTestFile(t, 2*DefaultBlockSize, copy)
	err := DestructiveCat([]string{file}, NopCloser(output), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestMultiFile(t *testing.T) {
	output := NewCountingHashWriter()
	copy := NewCountingHashWriter()
	file1 := makeRandomTestFile(t, DefaultBlockSize, copy)
	file2 := makeRandomTestFile(t, DefaultBlockSize, copy)
	err := DestructiveCat([]string{file1, file2}, NopCloser(output), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestMultiBlockMultiFile(t *testing.T) {
	output := NewCountingHashWriter()
	copy := NewCountingHashWriter()
	file1 := makeRandomTestFile(t, 2*DefaultBlockSize, copy)
	file2 := makeRandomTestFile(t, 2*DefaultBlockSize, copy)
	err := DestructiveCat([]string{file1, file2}, NopCloser(output), DefaultBlockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestStrangeBlockSize(t *testing.T) {
	blockSize := int64(17)
	numBlocks := int64(1031)
	output := bytes.NewBuffer(nil)
	copy := bytes.NewBuffer(nil)
	file := makeRandomTestFile(t, blockSize*numBlocks, copy)
	err := DestructiveCat([]string{file}, NopCloser(output), blockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != copy.Len() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.Len(), copy.Len())
	}
	if bytes.Compare(output.Bytes(), copy.Bytes()) != 0 {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestHalfBlock(t *testing.T) {
	totalBytes := int64(1024 * 1024)
	blockSize := totalBytes * 2
	output := bytes.NewBuffer(nil)
	copy := bytes.NewBuffer(nil)
	file := makeRandomTestFile(t, totalBytes, copy)
	err := DestructiveCat([]string{file}, NopCloser(output), blockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != copy.Len() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.Len(), copy.Len())
	}
	if bytes.Compare(output.Bytes(), copy.Bytes()) != 0 {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestIncompleteLastBlock(t *testing.T) {
	blockSize := int64(1024 * 1024)
	totalBytes := 8*blockSize + blockSize/2
	output := bytes.NewBuffer(nil)
	copy := bytes.NewBuffer(nil)
	file := makeRandomTestFile(t, totalBytes, copy)
	err := DestructiveCat([]string{file}, NopCloser(output), blockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != copy.Len() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.Len(), copy.Len())
	}
	if bytes.Compare(output.Bytes(), copy.Bytes()) != 0 {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestMultipleSizes(t *testing.T) {
	blockSize := int64(4096)
	sizes := []int64{1117, 15131, 28979, 8167, 9719, 10151, 14057, 29017, 17581, 2243, 36767, 29947, 28879, 15671, 17749, 27947}
	output := NewCountingHashWriter()
	copy := NewCountingHashWriter()
	files := make([]string, 0, len(sizes))
	for _, size := range sizes {
		files = append(files, makeRandomTestFile(t, size, copy))
	}
	err := DestructiveCat(files, NopCloser(output), blockSize)
	assertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func assertNoErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("%s: %v", msg, err)
		t.FailNow()
	}
}

func makeRandomTestFile(t *testing.T, size int64, copyTo io.Writer) string {
	file, err := ioutil.TempFile("", "destructive_cat_test.*.data")
	assertNoErr(t, err, "cannot create test file")
	name := file.Name()
	t.Cleanup(func() {
		_ = os.Remove(name)
	})
	writer := io.Writer(file)
	if copyTo != nil {
		writer = io.MultiWriter(file, copyTo)
	}
	_, err = io.CopyN(writer, rand.Reader, size)
	assertNoErr(t, err, "cannot write test data")
	err = file.Close()
	assertNoErr(t, err, "cannot close test file")
	return name
}

type CountingHashWriter struct {
	hash         hash.Hash
	bytesWritten int
}

func (h *CountingHashWriter) Write(p []byte) (n int, err error) {
	n, err = h.hash.Write(p)
	h.bytesWritten += n
	return n, err
}

func (h *CountingHashWriter) BytesWritten() int {
	return h.bytesWritten
}

func (a *CountingHashWriter) Equals(b *CountingHashWriter) bool {
	return bytes.Compare(a.hash.Sum(nil), b.hash.Sum(nil)) == 0
}

func NewCountingHashWriter() *CountingHashWriter {
	return &CountingHashWriter{
		hash:         sha512.New(),
		bytesWritten: 0,
	}
}

// Custom NopCloser returning an io.WriteCloser as io.Util only supports io.ReadCloser
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func NopCloser(r io.Writer) io.WriteCloser {
	return nopCloser{r}
}
