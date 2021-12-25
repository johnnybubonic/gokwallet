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
		It's mostly a placeholder for more specific errors.
	*/
	ErrNotInitialized error = errors.New("object not properly initialized")
	/*
		ErrOperationFailed is a generic failure message that will occur of a Dbus operation returns non-success.
		It is a placeholder for more specific messages.
	*/
	ErrOperationFailed error = errors.New("a Dbus operation has failed to execute successfully")
	/*
		ErrNoCreate is triggered if attempting to create an item (Folder, Password, etc.) but it fails.
		It is a placeholder for more specific messages.
	*/
	ErrNoCreate error = errors.New("failed to create an object")
	// ErrNoDisconnect can occur if trying to disconnect a Wallet from a WalletManager/application and a failure occurs.
	ErrNoDisconnect error = errors.New("failed to disconnect wallet from application")
	// ErrInvalidMap will get triggered if a populated map[string]string (even an empty one) is expected but a nil is received.
	ErrInvalidMap error = errors.New("invalid map; cannot be nil")
)

// Dbus Operation failures.
var (
	// ErrDbusOpfailNoHandle returns when attempting to open a Wallet and assign to Wallet.handle but received a nil handle.
	ErrDbusOpfailNoHandle error = errors.New("a wallet handler request returned nil")
	// ErrDbusOpfailRemoveFolder occurs when attempting to delete/remove a Folder from a Wallet but it did not complete successfully.
	ErrDbusOpfailRemoveFolder error = errors.New("failed to remove/delete a Folder from a Wallet")
)

// Initialization errors. They are more "detailed" ErrNotInitialized errors.
var (
	// ErrInitWM occurs if a WalletManager is not initialized properly.
	ErrInitWM error = errors.New("a WalletManager was not properly initialized")
	// ErrInitWallet occurs if a Wallet is not initialized properly.
	ErrInitWallet error = errors.New("a Wallet was not properly initialized")
	// ErrInitFolder occurs if a Folder is not initialized properly.
	ErrInitFolder error = errors.New("a Folder was not properly initialized")
	// ErrInitBlob occurs if a Blob is not initialized properly.
	ErrInitBlob error = errors.New("a Blob was not properly initialized")
	// ErrInitMap occurs if a Map is not initialized properly.
	ErrInitMap error = errors.New("a Map was not properly initialized")
	// ErrInitPassword occurs if a Password is not initialized properly.
	ErrInitPassword error = errors.New("a Password was not properly initialized")
	// ErrInitUnknownItem occurs if an UnknownItem is not initialized properly.
	ErrInitUnknownItem error = errors.New("an UnknownItem was not properly initialized")
)
