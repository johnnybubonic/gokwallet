package gokwallet

import (
	"bytes"
	"testing"
)

// TestUnknownItem tests all functions of an UnknownItem.
func TestUnknownItem(t *testing.T) {

	var u *UnknownItem
	var e *testEnv
	var err error
	var i int

	if e, err = getTestEnv(t); err != nil {
		t.Fatalf("failure getting test env: %v", err)
	}
	defer e.cleanup(t)

	if u, err = NewUnknownItem(e.f, unknownItemTest.String(), e.r); err != nil {
		t.Fatalf("failure getting UnknownItem: %v", err)
	}

	t.Logf("created UnknownItem '%v:%v:%v'; initialized: %v", e.w.Name, e.f.Name, u.Name, u.isInit)

	if err = u.SetValue(testBytes); err != nil {
		t.Errorf("failed to set value for UnknownItem '%v:%v:%v': %v", e.w.Name, e.f.Name, u.Name, err)
	}

	if i = bytes.Compare(u.Value, testBytes); i != 0 {
		t.Errorf("value '%#v' does not match expected value '%#v'", u.Value, testBytes)
	} else {
		t.Logf("value for UnknownItem '%v:%v:%v': %#v", e.w.Name, e.f.Name, u.Name, u.Value)
	}

	if err = u.SetValue(testBytesReplace); err != nil {
		t.Errorf("failed to set replacement value for UnknownItem '%v:%v:%v': %v", e.w.Name, e.f.Name, u.Name, err)
	}

	if err = u.Update(); err != nil {
		t.Errorf("failed to update UnknownItem '%v:%v:%v': %v", e.w.Name, e.f.Name, u.Name, err)
	}

	t.Logf("replacement value for UnknownItem '%v:%v:%v': %v", e.w.Name, e.f.Name, u.Name, u.Value)

	if i = bytes.Compare(u.Value, testBytesReplace); i != 0 {
		t.Errorf("value '%#v' does not match expected value '%#v'", u.Value, testBytesReplace)
	}

	if err = u.Delete(); err != nil {
		t.Errorf("failed to delete UnknownItem '%v:%v:%v': %v", e.w.Name, e.f.Name, u.Name, err)
	}

}
