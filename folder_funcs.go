package gokwallet

/*
	NewF returns a Wallet. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a WalletManager and wallet name.
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
		// handle:     0,
		isInit: false,
	}

	if err = folder.folderCheck(); err != nil {
		return
	}

	if folder.Recurse.All || folder.Recurse.Wallets {
		if err = folder.Update(); err != nil {
			return
		}
	}

	folder.isInit = true

	return
}

func (f *Folder) Update() (err error) {

	// TODO.

	return
}

func (f *Folder) folderCheck() (err error) {

	// TODO.

	return
}
