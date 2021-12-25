package gokwallet

import (
	"testing"
)

func TestNewBlob(t *testing.T) {

	var err error
	var f *Folder
	var b *Blob

	if _, _, f, err = getTestEnv(t); err != nil {
		t.Fatalf("failure getting test env: %v", err)
	}

	// if b, err = NewBlob(f)
	_ = f
	_ = b
}
