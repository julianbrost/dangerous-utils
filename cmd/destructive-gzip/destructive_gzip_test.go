package main

import (
	"compress/gzip"
	"dangerous-utils/internal/testutils"
	"io"
	"os"
	"syscall"
	"testing"
)

func TestSingle(t *testing.T) {
	size := int64(8 * 1024 * 1024)
	copy := testutils.NewCountingHashWriter()
	src := testutils.MakeRandomTestFile(t, size, copy)
	destFile := testutils.MakeEmptyTestFile(t)
	args := Args{
		Source:      src,
		Destination: destFile.Name(),
		Final:       true,
	}
	err := DestructiveGzip(&args)
	testutils.AssertNoErr(t, err, "DestructiveGzip returned an error")
	result := gzipFileToHash(t, destFile)
	testutils.AssertHashesEq(t, result, copy)
}

func TestMulti(t *testing.T) {
	calls := 4
	blockSize := int64(1024 * 1024)
	totalSize := int64(calls-1)*blockSize + 1 // ensure incomplete block at the end
	copy := testutils.NewCountingHashWriter()
	src := testutils.MakeRandomTestFile(t, totalSize, copy)
	destFile := testutils.MakeEmptyTestFile(t)
	for i := 0; i < calls; i++ {
		args := Args{
			Source:      src,
			Destination: destFile.Name(),
			Final:       i == calls-1,
			Max:         uint64(blockSize),
		}
		err := DestructiveGzip(&args)
		testutils.AssertNoErr(t, err, "DestructiveGzip returned an error")
	}
	result := gzipFileToHash(t, destFile)
	testutils.AssertHashesEq(t, result, copy)
}

func TestPartialDelete(t *testing.T) {
	calls := 4
	blockSize := getFsBlockSize(t)
	totalSize := int64(calls) * blockSize
	copy := testutils.NewCountingHashWriter()
	src := testutils.MakeRandomTestFile(t, totalSize, copy)
	destFile := testutils.MakeEmptyTestFile(t)
	for i := 0; i < calls; i++ {
		var stat1, stat2 syscall.Stat_t
		err := syscall.Stat(src, &stat1)
		testutils.AssertNoErr(t, err, "stat on source file")
		args := Args{
			Source:      src,
			Destination: destFile.Name(),
			Final:       i == calls-1,
			Max:         uint64(blockSize),
		}
		err = DestructiveGzip(&args)
		testutils.AssertNoErr(t, err, "DestructiveGzip returned an error")
		err = syscall.Stat(src, &stat2)
		testutils.AssertNoErr(t, err, "stat on source file")
		if stat2.Blocks == stat1.Blocks-1 {
			t.Error("file did not shrink by one block")
		}
	}
	result := gzipFileToHash(t, destFile)
	testutils.AssertHashesEq(t, result, copy)
}

func TestStrangeSizes(t *testing.T) {
	sizes := []int64{39727, 32887, 21929, 33091, 33049, 39847, 20563, 13711, 13691, 10193, 6547, 9551, 15601, 5261, 28309, 8821}
	totalSize := int64(0)
	for _, size := range sizes {
		totalSize += size
	}
	copy := testutils.NewCountingHashWriter()
	src := testutils.MakeRandomTestFile(t, totalSize, copy)
	destFile := testutils.MakeEmptyTestFile(t)
	for i, size := range sizes {
		args := Args{
			Source:      src,
			Destination: destFile.Name(),
			Final:       i == len(sizes)-1,
			Max:         uint64(size),
		}
		err := DestructiveGzip(&args)
		testutils.AssertNoErr(t, err, "DestructiveGzip returned an error")
	}
	result := gzipFileToHash(t, destFile)
	testutils.AssertHashesEq(t, result, copy)
}

func gzipFileToHash(t *testing.T, file *os.File) *testutils.CountingHashWriter {
	gzipReader, err := gzip.NewReader(file)
	testutils.AssertNoErr(t, err, "gzip.NewReader returned an error")
	hash := testutils.NewCountingHashWriter()
	_, err = io.Copy(hash, gzipReader)
	testutils.AssertNoErr(t, err, "reading gzip returned an error")
	return hash
}

func getFsBlockSize(t *testing.T) int64 {
	probeFile := testutils.MakeEmptyTestFile(t)
	var probeStat syscall.Stat_t
	err := syscall.Stat(probeFile.Name(), &probeStat)
	_ = probeFile.Close()
	testutils.AssertNoErr(t, err, "stat on probe file")
	return probeStat.Blksize
}
