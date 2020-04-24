package testutils

import "testing"

func AssertNoErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("%s: %v", msg, err)
		t.FailNow()
	}
}
