package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewFolder returns a Folder. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Wallet.
*/
func NewFolder(w *Wallet, name string, recursion *RecurseOpts) (folder *Folder, err error) {

	if !w.isInit {
		err = ErrNotInitialized
		return
	}

	folder = &Folder{
		DbusObject: w.DbusObject,
		Name:       name,
		Passwords:  nil,
		Maps:       nil,
		BinaryData: nil,
		Unknown:    nil,
		Recurse:    recursion,
		wm:         w.wm,
		wallet:     w,
		isInit:     false,
	}

	if folder.Recurse.AllWalletItems ||
		folder.Recurse.Passwords ||
		folder.Recurse.Maps ||
		folder.Recurse.Blobs ||
		folder.Recurse.UnknownItems {
		if err = folder.Update(); err != nil {
			return
		}
	}

	folder.isInit = true

	return
}

// HasEntry specifies if a Folder has an entry (WalletItem item) by the give entryName.
func (f *Folder) HasEntry(entryName string) (hasEntry bool, err error) {

	if err = f.Dbus.Call(
		DbusWMHasEntry, 0, f.wallet.handle, f.Name, entryName, f.wallet.wm.AppID,
	).Store(&hasEntry); err != nil {
		return
	}

	return
}

/*
	KeyNotExist returns true if a key/entry name entryName does *not* exist.
	Essentially the same as Folder.HasEntry, but whereas Folder.HasEntry requires the parent wallet
	to be open/unlocked, Folder.KeyNotExist does not require this.
*/
func (f *Folder) KeyNotExist(entryName string) (doesNotExist bool, err error) {

	if err = f.Dbus.Call(
		DbusWMKeyNotExist, 0, f.wallet.Name, f.Name, entryName,
	).Store(&doesNotExist); err != nil {
		return
	}

	return
}

// ListEntries lists all entries (WalletItem items) in a Folder (regardless of type) by name.
func (f *Folder) ListEntries() (entryNames []string, err error) {

	if err = f.Dbus.Call(
		DbusWMEntryList, 0, f.wallet.handle, f.Name, f.wallet.wm.AppID,
	).Store(&entryNames); err != nil {
		return
	}

	return
}

// RemoveEntry removes a WalletItem from a Folder given its entryName (key).
func (f *Folder) RemoveEntry(entryName string) (err error) {

	var rslt int32

	if err = f.Dbus.Call(
		DbusWMRemoveEntry, 0, f.wallet.handle, f.Name, entryName, f.wallet.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// RenameEntry renames a WalletItem in a Folder from entryName to newEntryName.
func (f *Folder) RenameEntry(entryName, newEntryName string) (err error) {

	var rslt int32

	if err = f.Dbus.Call(
		DbusWMRenameEntry, 0, f.wallet.handle, f.Name, entryName, newEntryName, f.wallet.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// Update runs all of the configured Update[type] methods for a Folder, depending on Folder.Recurse configuration.
func (f *Folder) Update() (err error) {

	var errs []error = make([]error, 0)

	if f.Recurse.AllWalletItems || f.Recurse.Passwords {
		if err = f.UpdatePasswords(); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}
	if f.Recurse.AllWalletItems || f.Recurse.Maps {
		if err = f.UpdateMaps(); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}
	if f.Recurse.AllWalletItems || f.Recurse.Blobs {
		if err = f.UpdateBlobs(); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}
	if f.Recurse.AllWalletItems || f.Recurse.UnknownItems {
		if err = f.UpdateUnknowns(); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}

// UpdatePasswords updates (populates) a Folder's Folder.Passwords.
func (f *Folder) UpdatePasswords() (err error) {

	var mapKeys []string
	var variant dbus.Variant
	var errs []error = make([]error, 0)

	if !f.isInit {
		err = ErrNotInitialized
		return
	}

	if err = f.Dbus.Call(
		DbusWMPasswordList, 0, f.wallet.handle, f.Name, f.wallet.wm.AppID,
	).Store(&variant); err != nil {
		return
	}

	mapKeys = bytemapKeys(variant)

	f.Passwords = make(map[string]*Password, len(mapKeys))

	for _, k := range mapKeys {
		if f.Passwords[k], err = NewPassword(f, k, f.Recurse); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}

// UpdateMaps updates (populates) a Folder's Folder.Maps.
func (f *Folder) UpdateMaps() (err error) {

	var mapKeys []string
	var variant dbus.Variant
	var errs []error = make([]error, 0)

	if err = f.Dbus.Call(
		DbusWMMapList, 0, f.wallet.handle, f.Name, f.wallet.wm.AppID,
	).Store(&variant); err != nil {
		return
	}

	mapKeys = bytemapKeys(variant)

	f.Maps = make(map[string]*Map, len(mapKeys))

	for _, k := range mapKeys {
		if f.Maps[k], err = NewMap(f, k, f.Recurse); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}

// UpdateBlobs updates (populates) a Folder's Folder.BinaryData.
func (f *Folder) UpdateBlobs() (err error) {

	var mapKeys []string
	var isBlob bool
	var variant dbus.Variant
	var errs []error = make([]error, 0)

	if !f.isInit {
		err = ErrNotInitialized
		return
	}

	if err = f.Dbus.Call(
		DbusWMEntriesList, 0, f.wallet.handle, f.Name, f.wallet.wm.AppID,
	).Store(&variant); err != nil {
		return
	}

	mapKeys = bytemapKeys(variant)

	f.BinaryData = make(map[string]*Blob, len(mapKeys))

	for _, k := range mapKeys {
		if isBlob, err = f.isType(k, KwalletdEnumTypeStream); err != nil {
			errs = append(errs, err)
			err = nil
			continue
		}
		if !isBlob {
			continue
		}

		if f.BinaryData[k], err = NewBlob(f, k, f.Recurse); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}

// UpdateUnknowns updates (populates) a Folder's Folder.Unknown.
func (f *Folder) UpdateUnknowns() (err error) {

	var mapKeys []string
	var isUnknown bool
	var variant dbus.Variant
	var errs []error = make([]error, 0)

	if !f.isInit {
		err = ErrNotInitialized
		return
	}

	if err = f.Dbus.Call(
		DbusWMEntriesList, 0, f.wallet.handle, f.Name, f.wallet.wm.AppID,
	).Store(&variant); err != nil {
		return
	}

	mapKeys = bytemapKeys(variant)

	f.Unknown = make(map[string]*UnknownItem, len(mapKeys))

	for _, k := range mapKeys {
		if isUnknown, err = f.isType(k, KwalletdEnumTypeUnknown); err != nil {
			errs = append(errs, err)
			err = nil
			continue
		}
		if !isUnknown {
			continue
		}

		if f.Unknown[k], err = NewUnknownItem(f, k, f.Recurse); err != nil {
			errs = append(errs, err)
			err = nil
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}

// WriteBlob adds or replaces a Blob to/in a Folder.
func (f *Folder) WriteBlob(entryName string, entryValue []byte) (err error) {

	if err = f.WriteEntry(entryName, KwalletdEnumTypeStream, entryValue); err != nil {
		return
	}

	return
}

/*
	WriteEntry is used for adding/replacing a WalletItem as a general interface.
	If possible, you'll want to use a item-type-specific method (e.g. Folder.WritePassword) as this one is a little unwieldy to use.
	entryType must be the relevant KwalletdEnumType* constant (do not use KwalletdEnumTypeUnused).
*/
func (f *Folder) WriteEntry(entryName string, entryType kwalletdEnumType, entryValue []byte) (err error) {

	var rslt int32

	if entryType == KwalletdEnumTypeUnused {
		err = ErrNoCreate
		return
	}

	if err = f.Dbus.Call(
		DbusWMWriteEntry, 0, f.wallet.handle, f.Name, entryName, entryValue, int32(entryType), f.wallet.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// WriteMap adds or replaces a Map to/in a Folder.
func (f *Folder) WriteMap(entryName string, entryValue map[string]string) (err error) {

	var rslt int32
	var b []byte

	if b, err = mapToBytes(entryValue); err != nil {
		return
	}

	if err = f.Dbus.Call(
		DbusWMWriteMap, 0, f.wallet.handle, f.Name, entryName, b, f.wallet.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// WritePassword adds or replaces a Password to/in a Folder.
func (f *Folder) WritePassword(entryName, entryValue string) (err error) {

	var rslt int32

	if err = f.Dbus.Call(
		DbusWMWritePassword, 0, f.wallet.handle, f.Name, entryName, entryValue, f.wallet.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// WriteUnknown adds or replaces an UnknownItem to/in a Folder.
func (f *Folder) WriteUnknown(entryName string, entryValue []byte) (err error) {

	if err = f.WriteEntry(entryName, KwalletdEnumTypeUnknown, entryValue); err != nil {
		return
	}

	return
}

// isType checks if a certain key keyName is of type typeCheck (via KwalletdEnumType*).
func (f *Folder) isType(keyName string, typeCheck kwalletdEnumType) (isOfType bool, err error) {

	var entryType int32

	if err = f.Dbus.Call(
		DbusWMEntryType, 0, f.wallet.handle, f.Name, keyName, f.wallet.wm.AppID,
	).Store(&entryType); err != nil {
		return
	}

	if int32(typeCheck) == entryType {
		isOfType = true
	}

	return
}
