package testutils

import "testing"

func AssertNoErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("%s: %v", msg, err)
		t.FailNow()
	}
}

func AssertHashesEq(t *testing.T, result, expected *CountingHashWriter) {
	if result.BytesWritten() != expected.BytesWritten() {
		t.Errorf("got %d bytes, expected %d", result.BytesWritten(), expected.BytesWritten())
	}
	if !result.Equals(expected) {
		t.Errorf("data did not match")
	}
}
