package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewBlob returns a Blob. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewBlob(f *Folder, keyName string, recursion *RecurseOpts) (blob *Blob, err error) {

	if !f.isInit {
		err = ErrInitFolder
		return
	}

	blob = &Blob{
		DbusObject: f.DbusObject,
		Name:       keyName,
		// Value:      "",
		Recurse: recursion,
		wm:      f.wallet.wm,
		wallet:  f.wallet,
		folder:  f,
		isInit:  false,
	}
	blob.isInit = true

	if blob.Recurse.AllWalletItems || blob.Recurse.Blobs {
		if err = blob.Update(); err != nil {
			return
		}
	}

	blob.isInit = true

	return
}

// Delete will delete this Blob from its parent Folder. You may want to run Folder.UpdateBlobs to update the existing map of Blob items.
func (b *Blob) Delete() (err error) {

	if err = b.folder.RemoveEntry(b.Name); err != nil {
		return
	}

	b = nil

	return
}

// SetValue will replace this Blob's Blob.Value.
func (b *Blob) SetValue(newValue []byte) (err error) {

	if _, err = b.folder.WriteBlob(b.Name, newValue); err != nil {
		return
	}

	b.Value = newValue

	return
}

// Update fetches a Blob's Blob.Value.
func (b *Blob) Update() (err error) {

	var call *dbus.Call
	var v dbus.Variant

	if call = b.Dbus.Call(
		DbusWMReadEntry, 0, b.folder.wallet.handle, b.folder.Name, b.Name, b.folder.wallet.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&v); err != nil {
		return
	}

	b.Value = v.Value().([]byte)

	return
}

// isWalletItem is needed for interface membership.
func (b *Blob) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
