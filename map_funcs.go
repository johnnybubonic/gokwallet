package gokwallet

/*
	NewMap returns a Map. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewMap(f *Folder, keyName string, recursion *RecurseOpts) (m *Map, err error) {

	if !f.isInit {
		err = ErrNotInitialized
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

	if m.Recurse.AllWalletItems || m.Recurse.Maps {
		if err = m.Update(); err != nil {
			return
		}
	}

	m.isInit = true

	return
}

// Update fetches a Map's Map.Value.
func (m *Map) Update() (err error) {

	var b []byte

	if err = m.Dbus.Call(
		DbusWMReadMap, 0, m.folder.wallet.handle, m.folder.Name, m.Name, m.folder.wallet.wm.AppID,
	).Store(&b); err != nil {
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
