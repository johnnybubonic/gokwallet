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
		if isBlob, err = f.isType(k, kwalletdEnumTypeStream); err != nil {
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
		if isUnknown, err = f.isType(k, kwalletdEnumTypeUnknown); err != nil {
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

// isType checks if a certain key keyName is of type typeCheck (via kwalletdEnumType*).
func (f *Folder) isType(keyName string, typeCheck int32) (isOfType bool, err error) {

	var entryType int32

	if err = f.Dbus.Call(
		DbusWMEntryType, 0, f.wallet.handle, f.Name, keyName, f.wallet.wm.AppID,
	).Store(&entryType); err != nil {
		return
	}

	return
}
