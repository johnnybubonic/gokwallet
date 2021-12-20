package gokwallet

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

// Update fetches an UnknownItem's UnknownItem.Value.
func (u *UnknownItem) Update() (err error) {

	// TODO.

	return
}

// isWalletItem is needed for interface membership.
func (u *UnknownItem) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
