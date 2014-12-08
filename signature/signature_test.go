package signature

import (
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	signature := Sign("abcdef", "x=1&y=2&z=3")

	if signature != "c9516346abf62876b6345817dba2f9a0c797ef26" {
		t.Errorf("Signature fails: %s", signature)
	}
}

func TestAppendSign(t *testing.T) {
	qs := AppendSign("abcdef", "x=1&y=2&z=3")

	if qs != "x=1&y=2&z=3&sig=c9516346abf62876b6345817dba2f9a0c797ef26" {
		t.Errorf("Append fails: %s", qs)
	}
}

func TestVerifySign(t *testing.T) {
	sign := "c9516346abf62876b6345817dba2f9a0c797ef26"

	same := VerifySign("abcdef", fmt.Sprintf("x=1&y=2&z=3&sig=%s", sign))

	if same == false {
		t.Errorf("Signature should be found in query string")
	}
}
