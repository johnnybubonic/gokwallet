package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewPassword returns a Password. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewPassword(f *Folder, keyName string, recursion *RecurseOpts) (password *Password, err error) {

	if !f.isInit {
		err = ErrInitFolder
		return
	}

	password = &Password{
		DbusObject: f.DbusObject,
		Name:       keyName,
		// Value:      "",
		Recurse: recursion,
		wm:      f.wallet.wm,
		wallet:  f.wallet,
		folder:  f,
		isInit:  false,
	}

	password.isInit = true

	if password.Recurse.AllWalletItems || password.Recurse.Passwords {
		if err = password.Update(); err != nil {
			return
		}
	}

	password.isInit = true

	return
}

// Delete will delete this Password from its parent Folder. You may want to run Folder.UpdatePasswords to update the existing map of Password items.
func (p *Password) Delete() (err error) {

	if err = p.folder.RemoveEntry(p.Name); err != nil {
		return
	}

	p = nil

	return
}

// SetValue will replace this Password's Password.Value.
func (p *Password) SetValue(newValue string) (err error) {

	if _, err = p.folder.WritePassword(p.Name, newValue); err != nil {
		return
	}

	p.Value = newValue

	return
}

// Update fetches a Password's Password.Value.
func (p *Password) Update() (err error) {

	var call *dbus.Call
	var b []byte

	if call = p.Dbus.Call(
		DbusWMReadPassword, 0, p.folder.wallet.handle, p.folder.Name, p.Name, p.folder.wallet.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&b); err != nil {
		return
	}

	p.Value = string(b)

	return
}

// isWalletItem is needed for interface membership.
func (p *Password) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
