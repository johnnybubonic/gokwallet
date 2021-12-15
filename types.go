package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	MultiError is a type of error.Error that can contain multiple error.Errors. Confused? Don't worry about it.
*/
type MultiError struct {
	// Errors is a slice of errors to combine/concatenate when .Error() is called.
	Errors []error `json:"errors"`
	// ErrorSep is a string to use to separate errors for .Error(). The default is "\n".
	ErrorSep string `json:"separator"`
}

// ConnPathCheckResult contains the result of validConnPath.
type ConnPathCheckResult struct {
	// ConnOK is true if the dbus.Conn is valid.
	ConnOK bool `json:"conn"`
	// PathOK is true if the Dbus path given is a valid type and value.
	PathOK bool `json:"path"`
}

// DbusObject is a base struct type to be anonymized by other types.
type DbusObject struct {
	// Conn is an active connection to the Dbus.
	Conn *dbus.Conn `json:"-"`
	// Dbus is the Dbus bus object.
	Dbus dbus.BusObject `json:"-"`
}

/*
	WalletManager is a general KWallet interface, sort of a handler for Dbus.
	It's used for fetching Wallet objects.
*/
type WalletManager struct {
	*DbusObject
	/*
		AppID is the application ID.
		The default is DefaultAppID.
	*/
	AppID string `json:"app_id"`
	/*
		Wallets is the collection of Wallets accessible in/to this WalletManager.
		Wallet.Name is the map key.
	*/
	Wallets map[string]*Wallet `json:"wallets"`
}

// Wallet contains one or more (or none) Folder objects.
type Wallet struct {
	*DbusObject
	// Name is the name of this Wallet.
	Name string `json:"name"`
	/*
		Folders contains all Folder objects in this Wallet.
		Folder.Name is the map key.
	*/
	Folders map[string]*Folder `json:"folders"`
}

// Folder contains secret object collections of Password, Map, Blob, and UnknownItem objects.
type Folder struct {
	*DbusObject
	// Name is the name of this Folder.
	Name string `json:"name"`
	/*
		Passwords contains a map of all Password objects in this Folder.
		Password.Name is the map key.
	*/
	Passwords map[string]*Password `json:"passwords"`
	/*
		Maps contains a map of all Map objects in this Folder.
		Map.Name is the map key.
	*/
	Maps map[string]*Map `json:"maps"`
	/*
		BinaryData contains a map if all Blob objects in this Folder.
		Blob.Name is the map key.
	*/
	BinaryData map[string]*Blob `json:"binary_data"`
	/*
		Unknown contains a map of all UnknownItem objects in this Folder.
		Unknown.Name is the map key.
	*/
	Unknown map[string]*UnknownItem `json:"unknown"`
}

// Password is a straightforward single-value secret of text.
type Password struct {
	*DbusObject
	// Name is the name of this Password.
	Name string `json:"name"`
	// Value is this Password's value.
	Value string `json:"value"`
}

// Map is a dictionary or key/value secret.
type Map struct {
	*DbusObject
	// Name is the name of this Map.
	Name string `json:"name"`
	// Value is this Map's value.
	Value map[string]string `json:"value"`
}

// Blob (binary large object, typographically BLOB) is secret binary data.
type Blob struct {
	*DbusObject
	// Name is the name of this Blob.
	Name string `json:"name"`
	// Value is this Blob's value.
	Value []byte `json:"value"`
}

/*
	UnknownItem is a secret item of unknown classification, so there isn't exactly a good way of determining a type for UnknownItem.Value.
	As such, its dbus.ObjectPath is used.
	TODO: There may be a method to fetch the raw bytes of the object (such as one would use with Blob) in the future.
*/
type UnknownItem struct {
	*DbusObject
	// Name is the name of this UnknownItem.
	Name string `json:"name"`
	// Value is the Dbus path of this UnknownItem.
	Value dbus.ObjectPath `json:"value"`
}

// WalletItem is an interface to manage wallet objects: Password, Map, Blob, or UnknownItem.
type WalletItem interface {
	isWalletItem() (isWalletItem bool)
}
