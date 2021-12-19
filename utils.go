package gokwallet

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
	}

	return
}
