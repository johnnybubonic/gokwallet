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
		(TODO: When wallet file support is added, the *filename* will be the map key.
			   This is to mitigate namespace conflicts between Dbus and file wallets.)
	*/
	Wallets map[string]*Wallet `json:"wallets"`
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// Enabled is true if KWalletD is enabled/running.
	Enabled bool `json:"enabled"`
	// Local is the "local" wallet.
	Local *Wallet `json:"local_wallet"`
	// Network is the "network" wallet.
	Network *Wallet `json:"network_wallet"`
	// isInit flags whether this is "properly" set up (i.e. was initialized via NewWalletManager).
	isInit bool
	// walletFiles are (resolved and vetted) wallet files (kwl, xml).
	walletFiles []string
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
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// IsUnlocked specifies if this Wallet is open ("unlocked") or not.
	IsUnlocked bool `json:"open"`
	/*
		FilePath is:
		- empty if this is an internal Wallet, or
		- the filepath to the wallet file if this is an on-disk wallet (either .kwl or .xml)
	*/
	FilePath string `json:"wallet_file"`
	// wm is the parent WalletManager this Wallet was fetched from.
	wm *WalletManager
	// handle is this Wallet's handler number.
	handle int32
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
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
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// wm is the parent WalletManager that Folder.wallet was fetched from.
	wm *WalletManager
	// wallet is the parent Wallet this Folder was fetched from.
	wallet *Wallet
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
}

// Password is a straightforward single-value secret of text.
type Password struct {
	*DbusObject
	// Name is the name of this Password.
	Name string `json:"name"`
	// Value is this Password's value.
	Value string `json:"value"`
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// wm is the parent WalletManager that Password.folder.wallet was fetched from.
	wm *WalletManager
	// wallet is the parent Wallet that Password.folder was fetched from.
	wallet *Wallet
	// folder is the parent Folder this Password was fetched from.
	folder *Folder
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
}

// Map is a dictionary or key/value secret.
type Map struct {
	*DbusObject
	// Name is the name of this Map.
	Name string `json:"name"`
	// Value is this Map's value.
	Value map[string]string `json:"value"`
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// wm is the parent WalletManager that Map.folder.wallet was fetched from.
	wm *WalletManager
	// wallet is the parent Wallet that Map.folder was fetched from.
	wallet *Wallet
	// folder is the parent Folder this Map was fetched from.
	folder *Folder
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
}

// Blob (binary large object, typographically BLOB) is secret binary data.
type Blob struct {
	*DbusObject
	// Name is the name of this Blob.
	Name string `json:"name"`
	// Value is this Blob's value.
	Value []byte `json:"value"`
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// wm is the parent WalletManager that Blob.folder.wallet was fetched from.
	wm *WalletManager
	// wallet is the parent Wallet that Blob.folder was fetched from.
	wallet *Wallet
	// folder is the parent Folder this Blob was fetched from.
	folder *Folder
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
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
	// Recurse contains the relevant RecurseOpts.
	Recurse *RecurseOpts `json:"recurse_opts"`
	// wm is the parent WalletManager that UnknownItem.folder.wallet was fetched from.
	wm *WalletManager
	// wallet is the parent Wallet that UnknownItem.folder was fetched from.
	wallet *Wallet
	// folder is the parent Folder this UnknownItem was fetched from.
	folder *Folder
	// isInit flags whether this is "properly" set up (i.e. has a handle).
	isInit bool
}

// WalletItem is an interface to manage wallet objects: Password, Map, Blob, or UnknownItem.
type WalletItem interface {
	isWalletItem() (isWalletItem bool)
}

/*
	RecurseOpts controls whether recursion should be done on objects when fetching them.
	E.g. if fetching a WalletManager (via NewWalletManager) and RecurseOpts.Wallet is true,
	then WalletManager.Wallets will be populated with Wallet objects.
*/
type RecurseOpts struct {
	/*
		All, if true, specifies that all possible recursions should be done.
		If true, it takes precedent over all over RecurseOpts fields (with the exception of RecurseOpts.AllWalletItems).

		Performed in/from:
		WalletManager
		Wallet
		Folder
		(WalletItem)
	*/
	All bool `json:"none"`
	/*
		Wallets, if true, indicates that Wallet objects should have Wallet.Update called.

		Performed in/from: WalletManager
	*/
	Wallets bool `json:"wallet"`
	/*
		Folders, if true, indicates that Folder objects should have Folder.Update called.

		Performed in/from:
		Wallet

		May be performed in/from (depending on other fields):
		WalletManager
	*/
	Folders bool `json:"folder"`
	/*
		AllWalletItems, if true, indicates that all WalletItem entries should have (WalletItem).Update() called.
		If true, it takes precedent over all over relevant RecurseOpts fields for each WalletItem type
		(i.e. RecurseOpts.Passwords, RecurseOpts.Maps, RecurseOpts.Blobs, RecurseOpts.UnknownItems).

		Performed in/from:
		Folder

		May be performed in/from (depending on other fields):
		WalletManager
		Wallet
	*/
	AllWalletItems bool `json:"wallet_item"`
	/*
		Passwords, if true, indicates that Password objects should have Password.Update() called.

		Performed in/from:
		Folder

		May be performed in/from (depending on other fields):
		WalletManager
		Wallet
	*/
	Passwords bool `json:"password"`
	/*
		Maps, if true, indicates that Map objects should have Map.Update() called.

		Performed in/from:
		Folder

		May be performed in/from (depending on other fields):
		WalletManager
		Wallet
	*/
	Maps bool `json:"map"`
	/*
		Blobs, if true, indicates that Blob objects should have Blob.Update() called.

		Performed in/from:
		Folder

		May be performed in/from (depending on other fields):
		WalletManager
		Wallet
	*/
	Blobs bool `json:"blob"`
	/*
		UnknownItems indicates that UnknownItem objects should have UnknownItem.Update() called.

		Performed in/from:
		Folder

		May be performed in/from (depending on other fields):
		WalletManager
		Wallet
	*/
	UnknownItems bool `json:"unknown_item"`
}
