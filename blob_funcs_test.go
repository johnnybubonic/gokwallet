package gokwallet

import (
	"bytes"
	"testing"
)

// TestBlob tests all functions of a Blob.
func TestBlob(t *testing.T) {

	var b *Blob
	var e *testEnv
	var err error
	var i int

	if e, err = getTestEnv(t); err != nil {
		t.Fatalf("failure getting test env: %v", err)
	}
	defer e.cleanup(t)

	if b, err = NewBlob(e.f, blobTest.String(), e.r); err != nil {
		t.Fatalf("failure getting Blob: %v", err)
	}

	t.Logf("created Blob '%v:%v:%v'; initialized: %v", e.w.Name, e.f.Name, b.Name, b.isInit)

	if err = b.SetValue(testBytes); err != nil {
		t.Errorf("failed to set value for Blob '%v:%v:%v': %v", e.w.Name, e.f.Name, b.Name, err)
	}

	if i = bytes.Compare(b.Value, testBytes); i != 0 {
		t.Errorf("value '%#v' does not match expected value '%#v'", b.Value, testBytes)
	} else {
		t.Logf("value for Blob '%v:%v:%v': %#v", e.w.Name, e.f.Name, b.Name, b.Value)
	}

	if err = b.SetValue(testBytesReplace); err != nil {
		t.Errorf("failed to set replacement value for Blob '%v:%v:%v': %v", e.w.Name, e.f.Name, b.Name, err)
	}

	if err = b.Update(); err != nil {
		t.Errorf("failed to update Blob '%v:%v:%v': %v", e.w.Name, e.f.Name, b.Name, err)
	}

	t.Logf("replacement value for Blob '%v:%v:%v': %v", e.w.Name, e.f.Name, b.Name, b.Value)

	if i = bytes.Compare(b.Value, testBytesReplace); i != 0 {
		t.Errorf("value '%#v' does not match expected value '%#v'", b.Value, testBytesReplace)
	}

	if err = b.Delete(); err != nil {
		t.Errorf("failed to delete Blob '%v:%v:%v': %v", e.w.Name, e.f.Name, b.Name, err)
	}

}
