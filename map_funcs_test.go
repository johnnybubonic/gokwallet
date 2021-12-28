package gokwallet

import (
	"reflect"
	"testing"
)

// TestMap tests all functions of Map.
func TestMap(t *testing.T) {

	var m *Map
	var e *testEnv
	var b bool
	var err error

	if e, err = getTestEnv(t); err != nil {
		t.Fatalf("failure getting test env: %v", err)
	}
	defer e.cleanup(t)

	if m, err = NewMap(e.f, mapTest.String(), e.r); err != nil {
		t.Fatalf("failure getting Map: %v", err)
	}

	t.Logf("created Map '%v:%v:%v'; initialized: %v", e.w.Name, e.f.Name, m.Name, m.isInit)

	if err = m.SetValue(testMap); err != nil {
		t.Errorf("failed to set value for Map '%v:%v:%v': %v", e.w.Name, e.f.Name, m.Name, err)
	}

	if b = reflect.DeepEqual(m.Value, testMap); !b {
		t.Errorf("value '%#v' does not match expected value '%#v'", m.Value, testMap)
	} else {
		t.Logf("value for Map '%v:%v:%v': %#v", e.w.Name, e.f.Name, m.Name, m.Value)
	}

	if err = m.SetValue(testMapReplace); err != nil {
		t.Errorf("failed to set replacement value for Map '%v:%v:%v': %v", e.w.Name, e.f.Name, m.Name, err)
	}

	if err = m.Update(); err != nil {
		t.Errorf("failed to update Map '%v:%v:%v': %v", e.w.Name, e.f.Name, m.Name, err)
	}

	t.Logf("replacement value for Map '%v:%v:%v': %v", e.w.Name, e.f.Name, m.Name, m.Value)

	if b = reflect.DeepEqual(m.Value, testMapReplace); !b {
		t.Errorf("value '%#v' does not match expected value '%#v'", m.Value, testMapReplace)
	}

	if err = m.Delete(); err != nil {
		t.Errorf("failed to delete Map '%v:%v:%v': %v", e.w.Name, e.f.Name, m.Name, err)
	}

}
