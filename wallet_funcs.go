package gokwallet

/*
	NewWallet returns a Wallet. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a WalletManager and wallet name.
*/
func NewWallet(wm *WalletManager, name string, recursion *RecurseOpts) (wallet *Wallet, err error) {

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	wallet = &Wallet{
		DbusObject: wm.DbusObject,
		Name:       name,
		Folders:    nil,
		Recurse:    recursion,
		wm:         wm,
		// handle:     0,
		isInit: false,
	}

	// TODO: remove this and leave to caller, since it might use PamOpen instead? Fail back to it?
	if err = wallet.walletCheck(); err != nil {
		return
	}

	if wallet.Recurse.All || wallet.Recurse.Wallets {
		if err = wallet.Update(); err != nil {
			return
		}
	}

	wallet.isInit = true

	return
}

// Close closes a Wallet.
func (w *Wallet) Close() (err error) {

	var rslt int32

	if !w.isInit {
		err = ErrNotInitialized
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if err = w.Dbus.Call(
		DbusWMClose, 0, w.handle, false, w.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

/*
	ForceClose is like Close but will still close a Wallet even if currently in use by the WalletManager.
	(Despite the insinuation by the name, this should be a relatively safe operation).
*/
func (w *Wallet) ForceClose() (err error) {

	var rslt int32

	if !w.isInit {
		err = ErrNotInitialized
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if err = w.Dbus.Call(
		DbusWMClose, 0, w.handle, true, w.wm.AppID,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return

	return
}

// IsOpen returns whether a Wallet is open ("unlocked") or not (as well as updates Wallet.IsOpen).
func (w *Wallet) IsOpen() (isOpen bool, err error) {

	// We can call the same method with w.handle instead of w.Name. We don't have a handler yet though.
	if err = w.Dbus.Call(
		DbusWMIsOpen, 0, w.Name,
	).Store(&w.IsUnlocked); err != nil {
		return
	}

	isOpen = w.IsUnlocked

	return
}

// ListFolders lists all Folder names in a Wallet.
func (w *Wallet) ListFolders() (folderList []string, err error) {

	if err = w.walletCheck(); err != nil {
		return
	}

	if err = w.Dbus.Call(
		DbusWMFolderList, 0, w.handle, w.wm.AppID,
	).Store(&folderList); err != nil {
		return
	}

	return
}

/*
	Open will open ("unlock") a Wallet.
	It will no-op if the Wallet is already open.
*/
func (w *Wallet) Open() (err error) {

	if _, err = w.IsOpen(); err != nil {
		return
	}

	if !w.IsUnlocked {
		if err = w.Dbus.Call(
			DbusWMOpen, 0,
		).Store(&w.handle); err != nil {
			return
		}
	}

	w.IsUnlocked = true

	return
}

// Update fetches/updates all Folder objects in a Wallet.
func (w *Wallet) Update() (err error) {

	var folderNames []string
	var errs []error = make([]error, 0)

	if folderNames, err = w.ListFolders(); err != nil {
		return
	}

	w.Folders = make(map[string]*Folder)

	for _, fn := range folderNames {
		if w.Folders[fn], err = NewFolder(w, fn, w.Recurse); err != nil {
			errs = append(errs, err)
			err = nil
			continue
		}
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
		return
	}

	return
}

// walletCheck will check if a Wallet is (initialized and) opened and, if not, attempt to open it.
func (w *Wallet) walletCheck() (err error) {

	if !w.isInit {
		err = ErrNotInitialized
		return
	}

	if _, err = w.IsOpen(); err != nil {
		return
	}

	if !w.IsUnlocked {
		if err = w.Open(); err != nil {
			return
		}
	}

	return
}
