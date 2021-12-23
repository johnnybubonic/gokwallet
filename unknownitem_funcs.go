package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewUnknownItem returns an UnknownItem. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewUnknownItem(f *Folder, keyName string, recursion *RecurseOpts) (unknown *UnknownItem, err error) {

	if !f.isInit {
		err = ErrNotInitialized
		return
	}

	unknown = &UnknownItem{
		DbusObject: f.DbusObject,
		Name:       keyName,
		// Value:      "",
		Recurse: recursion,
		wm:      f.wallet.wm,
		wallet:  f.wallet,
		folder:  f,
		isInit:  false,
	}

	if unknown.Recurse.AllWalletItems || unknown.Recurse.UnknownItems {
		if err = unknown.Update(); err != nil {
			return
		}
	}

	unknown.isInit = true

	return
}

// Delete will delete this UnknownItem from its parent Folder. You may want to run Folder.UpdateUnknowns to update the existing map of UnknownItem items.
func (u *UnknownItem) Delete() (err error) {

	if err = u.folder.RemoveEntry(u.Name); err != nil {
		return
	}

	u = nil

	return
}

// SetValue will replace this UnknownItem's UnknownItem.Value.
func (u *UnknownItem) SetValue(newValue []byte) (err error) {

	if _, err = u.folder.WriteUnknown(u.Name, newValue); err != nil {
		return
	}

	u.Value = newValue

	return
}

// Update fetches an UnknownItem's UnknownItem.Value.
func (u *UnknownItem) Update() (err error) {

	var v dbus.Variant

	if err = u.Dbus.Call(
		DbusWMReadEntry, 0, u.folder.wallet.handle, u.folder.Name, u.Name, u.folder.wallet.wm.AppID,
	).Store(&v); err != nil {
		return
	}

	u.Value = v.Value().([]byte)

	return
}

// isWalletItem is needed for interface membership.
func (u *UnknownItem) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
