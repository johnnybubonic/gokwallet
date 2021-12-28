package gokwallet

import (
	"testing"
)

// TestPassword tests all functions of a Password.
func TestPassword(t *testing.T) {

	var p *Password
	var e *testEnv
	var err error

	if e, err = getTestEnv(t); err != nil {
		t.Fatalf("failure getting test env: %v", err)
	}
	defer e.cleanup(t)

	if p, err = NewPassword(e.f, passwordTest.String(), e.r); err != nil {
		t.Fatalf("failure getting Password: %v", err)
	}

	t.Logf("created Password '%v:%v:%v'; initialized: %v", e.w.Name, e.f.Name, p.Name, p.isInit)

	if err = p.SetValue(testPassword); err != nil {
		t.Errorf("failed to set value for Password '%v:%v:%v': %v", e.w.Name, e.f.Name, p.Name, err)
	}

	if p.Value != testPassword {
		t.Errorf("value '%#v' does not match expected value '%#v'", p.Value, testPassword)
	} else {
		t.Logf("value for Password '%v:%v:%v': %#v", e.w.Name, e.f.Name, p.Name, p.Value)
	}

	if err = p.SetValue(testPasswordReplace); err != nil {
		t.Errorf("failed to set replacement value for Password '%v:%v:%v': %v", e.w.Name, e.f.Name, p.Name, err)
	}

	if err = p.Update(); err != nil {
		t.Errorf("failed to update Password '%v:%v:%v': %v", e.w.Name, e.f.Name, p.Name, err)
	}

	t.Logf("replacement value for Password '%v:%v:%v': %v", e.w.Name, e.f.Name, p.Name, p.Value)

	if p.Value != testPasswordReplace {
		t.Errorf("value '%#v' does not match expected value '%#v'", p.Value, testPasswordReplace)
	}

	if err = p.Delete(); err != nil {
		t.Errorf("failed to delete Password '%v:%v:%v': %v", e.w.Name, e.f.Name, p.Name, err)
	}

}
