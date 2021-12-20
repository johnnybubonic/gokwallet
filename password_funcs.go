package gokwallet

/*
	NewPassword returns a Password. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewPassword(f *Folder, keyName string, recursion *RecurseOpts) (password *Password, err error) {

	if !f.isInit {
		err = ErrNotInitialized
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

	if password.Recurse.AllWalletItems || password.Recurse.Passwords {
		if err = password.Update(); err != nil {
			return
		}
	}

	password.isInit = true

	return
}

// Update fetches a Password's Password.Value.
func (p *Password) Update() (err error) {

	// TODO.

	return
}

// isWalletItem is needed for interface membership.
func (p *Password) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
