package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	resultCheck checks the result code from a Dbus call and returns an error if not successful.
	See also resultPassed.
*/
func resultCheck(result int32) (err error) {

	// This is technically way more complex than it needs to be, but is extendable for future use.
	switch i := result; i {
	case DbusSuccess:
		err = nil
	case DbusFailure:
		err = ErrOperationFailed
	default:
		err = ErrOperationFailed
	}

	return
}

/*
	resultPassed checks the result code from a Dbus call and returns a boolean as to whether the result is pass or not.
	See also resultCheck.
*/
func resultPassed(result int32) (passed bool) {

	// This is technically way more complex than it needs to be, but is extendable for future use.
	switch i := result; i {
	case DbusSuccess:
		passed = true
	case DbusFailure:
		passed = false
	default:
		passed = false
	}

	return
}

// bytemapKeys is used to parse out Map names when fetching from Dbus.
func bytemapKeys(variant dbus.Variant) (keyNames []string) {

	var d map[string]dbus.Variant

	d = variant.Value().(map[string]dbus.Variant)

	keyNames = make([]string, len(d))

	idx := 0
	for k, _ := range d {
		keyNames[idx] = k
		idx++
	}

	return
}
