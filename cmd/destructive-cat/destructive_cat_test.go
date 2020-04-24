package main

import (
	"bytes"
	"dangerous-utils/internal/testutils"
	"io/ioutil"
	"os"
	"testing"
)

func TestEmptyFiles(t *testing.T) {
	output := bytes.NewBuffer(nil)
	err := DestructiveCat([]string{}, testutils.NopCloser(output), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != 0 {
		t.Errorf("DestructiveCat with an empty file list wrote %d bytes", output.Len())
	}
}

func TestSimple(t *testing.T) {
	output := bytes.NewBuffer(nil)
	copy := bytes.NewBuffer(nil)
	file := testutils.MakeRandomTestFile(t, DefaultBlockSize, copy)
	err := DestructiveCat([]string{file}, testutils.NopCloser(output), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	if output.Len() != copy.Len() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.Len(), copy.Len())
	}
	if bytes.Compare(output.Bytes(), copy.Bytes()) != 0 {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestInputDeleted(t *testing.T) {
	size := int64(1024)
	file := testutils.MakeRandomTestFile(t, size, nil)
	err := DestructiveCat([]string{file}, testutils.NopCloser(ioutil.Discard), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	_, err = os.Stat(file)
	if !os.IsNotExist(err) {
		t.Error("input file was not deleted")
	}
}

func TestMultiBlock(t *testing.T) {
	output := testutils.NewCountingHashWriter()
	copy := testutils.NewCountingHashWriter()
	file := testutils.MakeRandomTestFile(t, 2*DefaultBlockSize, copy)
	err := DestructiveCat([]string{file}, testutils.NopCloser(output), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestMultiFile(t *testing.T) {
	output := testutils.NewCountingHashWriter()
	copy := testutils.NewCountingHashWriter()
	file1 := testutils.MakeRandomTestFile(t, DefaultBlockSize, copy)
	file2 := testutils.MakeRandomTestFile(t, DefaultBlockSize, copy)
	err := DestructiveCat([]string{file1, file2}, testutils.NopCloser(output), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}

func TestMultiBlockMultiFile(t *testing.T) {
	output := testutils.NewCountingHashWriter()
	copy := testutils.NewCountingHashWriter()
	file1 := testutils.MakeRandomTestFile(t, 2*DefaultBlockSize, copy)
	file2 := testutils.MakeRandomTestFile(t, 2*DefaultBlockSize, copy)
	err := DestructiveCat([]string{file1, file2}, testutils.NopCloser(output), DefaultBlockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
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
	file := testutils.MakeRandomTestFile(t, blockSize*numBlocks, copy)
	err := DestructiveCat([]string{file}, testutils.NopCloser(output), blockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
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
	file := testutils.MakeRandomTestFile(t, totalBytes, copy)
	err := DestructiveCat([]string{file}, testutils.NopCloser(output), blockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
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
	file := testutils.MakeRandomTestFile(t, totalBytes, copy)
	err := DestructiveCat([]string{file}, testutils.NopCloser(output), blockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
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
	output := testutils.NewCountingHashWriter()
	copy := testutils.NewCountingHashWriter()
	files := make([]string, 0, len(sizes))
	for _, size := range sizes {
		files = append(files, testutils.MakeRandomTestFile(t, size, copy))
	}
	err := DestructiveCat(files, testutils.NopCloser(output), blockSize)
	testutils.AssertNoErr(t, err, "DestructiveCat returned an error")
	if output.BytesWritten() != copy.BytesWritten() {
		t.Errorf("DestructiveCat wrote %d bytes, expected %d", output.BytesWritten(), copy.BytesWritten())
	}
	if !output.Equals(copy) {
		t.Errorf("DestructiveCat did not write the expected data")
	}
}
