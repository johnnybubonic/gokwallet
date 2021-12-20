package gokwallet

import (
	"errors"
)

var (
	/*
		ErrNotInitialized will be triggered if attempting to interact with an object that has not been properly initialized.
		Notably, in most/all cases this means that it was not created via a New<object> func (for instance,
		this would lead to a Wallet missing a handler).
		It is intended as a safety check (so that you don't accidentally delete a wallet with e.g. a handler of 0 when
		trying to delete a different wallet).
	*/
	ErrNotInitialized error = errors.New("object not properly initialized")
	/*
		ErrOperationFailed is a generic failure message that will occur of a Dbus operation returns non-success.
	*/
	ErrOperationFailed error = errors.New("a Dbus operation has failed to execute successfully")
	/*
		ErrNoCreate is triggered if attempting to create an item (Folder, Password, etc.) but it fails.
	*/
	ErrNoCreate error = errors.New("failed to create an object")
	// ErrNoDisconnect can occur if trying to disconnect a Wallet from a WalletManager/application and a failure occurs.
	ErrNoDisconnect error = errors.New("failed to disconnect wallet from application")
)
