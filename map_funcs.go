package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewMap returns a Map. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewMap(f *Folder, keyName string, recursion *RecurseOpts) (m *Map, err error) {

	if !f.isInit {
		err = ErrInitFolder
		return
	}

	m = &Map{
		DbusObject: f.DbusObject,
		Name:       keyName,
		// Value:      "",
		Recurse: recursion,
		wm:      f.wallet.wm,
		wallet:  f.wallet,
		folder:  f,
		isInit:  false,
	}

	m.isInit = true

	if m.Recurse.AllWalletItems || m.Recurse.Maps {
		if err = m.Update(); err != nil {
			return
		}
	}

	m.isInit = true

	return
}

// Delete will delete this Map from its parent Folder. You may want to run Folder.UpdateMaps to update the existing map of Map items.
func (m *Map) Delete() (err error) {

	if err = m.folder.RemoveEntry(m.Name); err != nil {
		return
	}

	m = nil

	return
}

// Exists returns true if this Map actually exists.
func (m *Map) Exists() (exists bool, err error) {

	if exists, err = m.folder.HasEntry(m.Name); err != nil {
		return
	}

	return
}

// Rename renames this Map (changes its key).
func (m *Map) Rename(newName string) (err error) {

	if err = m.folder.RenameEntry(m.Name, newName); err != nil {
		return
	}

	m.Name = newName

	return
}

// SetValue will replace this Map's Map.Value.
func (m *Map) SetValue(newValue map[string]string) (err error) {

	if _, err = m.folder.WriteMap(m.Name, newValue); err != nil {
		return
	}

	m.Value = newValue

	return
}

// Update fetches a Map's Map.Value.
func (m *Map) Update() (err error) {

	var call *dbus.Call
	var b []byte

	if err = m.folder.wallet.walletCheck(); err != nil {
		return
	}

	if call = m.Dbus.Call(
		DbusWMReadMap, 0, m.folder.wallet.handle, m.folder.Name, m.Name, m.folder.wallet.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&b); err != nil {
		return
	}

	if m.Value, _, err = bytesToMap(b); err != nil {
		return
	}

	return
}

// isWalletItem is needed for interface membership.
func (m *Map) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
