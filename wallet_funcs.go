package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewWallet returns a Wallet. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a WalletManager and wallet name.
*/
func NewWallet(wm *WalletManager, name string, recursion *RecurseOpts) (wallet *Wallet, err error) {

	if !wm.isInit {
		err = ErrInitWM
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

	wallet.isInit = true

	// TODO: remove this and leave to caller, since it might use PamOpen instead? Fail back to it?
	if err = wallet.walletCheck(); err != nil {
		return
	}

	if wallet.Recurse.All || wallet.Recurse.Folders {
		if err = wallet.Update(); err != nil {
			return
		}
	}

	wallet.isInit = true

	return
}

// Disconnect disconnects this Wallet from its parent WalletManager.
func (w *Wallet) Disconnect() (err error) {

	var call *dbus.Call
	var ok bool

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMDisconnectApp, 0, w.Name, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&ok); err != nil {
		return
	}

	if !ok {
		err = ErrNoDisconnect
	}

	return
}

// DisconnectApplication disconnects this Wallet from a specified WalletManager/application (see Wallet.Connections).
func (w *Wallet) DisconnectApplication(appName string) (err error) {

	var call *dbus.Call
	var ok bool

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMDisconnectApp, 0, appName, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&ok); err != nil {
		return
	}

	if !ok {
		err = ErrNoDisconnect
	}

	return
}

/*
	ChangePassword will change (or set) the password for a Wallet.
	Note that this *must* be done via the windowing/graphical layer.
	There is no way to change a Wallet's password via the Dbus API.
*/
func (w *Wallet) ChangePassword() (err error) {

	var call *dbus.Call

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMChangePassword, 0, w.Name, DefaultWindowID, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}

	return
}

// Close closes a Wallet.
func (w *Wallet) Close() (err error) {

	var call *dbus.Call
	var rslt int32

	if err = w.walletCheck(); err != nil {
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if call = w.Dbus.Call(
		DbusWMClose, 0, w.handle, false, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// Connections lists the application names for connections to ("users of") this Wallet.
func (w *Wallet) Connections() (connList []string, err error) {

	var call *dbus.Call

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMUsers, 0, w.Name,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&connList); err != nil {
		return
	}

	return
}

// CreateFolder creates a new Folder in a Wallet.
func (w *Wallet) CreateFolder(name string) (err error) {

	var call *dbus.Call
	var ok bool

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMCreateFolder, 0, w.handle, name, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&ok); err != nil {
		return
	}

	if !ok {
		err = ErrNoCreate
	}

	return
}

// Delete deletes a Wallet.
func (w *Wallet) Delete() (err error) {

	var call *dbus.Call
	var rslt int32

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMDeleteWallet, 0, w.Name,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	w = nil

	return
}

// FolderExists indicates if a Folder exists in a Wallet or not.
func (w *Wallet) FolderExists(folderName string) (exists bool, err error) {

	var call *dbus.Call
	var notExists bool

	// We don't need a walletcheck here since we don't need a handle.

	if call = w.Dbus.Call(
		DbusWMFolderNotExist, 0, w.Name, folderName,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&notExists); err != nil {
		return
	}

	exists = !notExists

	return
}

/*
	ForceClose is like Close but will still close a Wallet even if currently in use by the WalletManager.
	(Despite the insinuation by the name, this should be a relatively safe operation).
*/
func (w *Wallet) ForceClose() (err error) {

	var call *dbus.Call
	var rslt int32

	if err = w.walletCheck(); err != nil {
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if call = w.Dbus.Call(
		DbusWMClose, 0, w.handle, true, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// HasFolder indicates if a Wallet has a Folder in it named folderName.
func (w *Wallet) HasFolder(folderName string) (hasFolder bool, err error) {

	var call *dbus.Call

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMHasFolder, 0, w.handle, folderName, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&hasFolder); err != nil {
		return
	}

	return
}

// IsOpen returns whether a Wallet is open ("unlocked") or not (as well as updates Wallet.IsOpen).
func (w *Wallet) IsOpen() (isOpen bool, err error) {

	var call *dbus.Call

	// We don't call walletcheck here because this method is called by a walletcheck.
	if !w.isInit {
		err = ErrInitWallet
		return
	}

	// We can call the same method with w.handle instead of w.Name. We don't have a handler yet though.
	if call = w.Dbus.Call(
		DbusWMIsOpen, 0, w.Name,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&w.IsUnlocked); err != nil {
		return
	}

	isOpen = w.IsUnlocked

	return
}

// ListFolders lists all Folder names in a Wallet.
func (w *Wallet) ListFolders() (folderList []string, err error) {

	var call *dbus.Call

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMFolderList, 0, w.handle, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&folderList); err != nil {
		return
	}

	return
}

/*
	Open will open ("unlock") a Wallet.
	It will no-op if the Wallet is already open.
*/
func (w *Wallet) Open() (err error) {

	var call *dbus.Call
	var handler *int32 = new(int32)

	// We don't call walletcheck here because this method is called by a walletcheck.
	if !w.isInit {
		err = ErrInitWallet
		return
	}

	if !w.IsUnlocked || !w.hasHandle {
		if call = w.Dbus.Call(
			DbusWMOpen, 0, w.Name, DefaultWindowID, w.wm.AppID,
		); call.Err != nil {
			err = call.Err
			return
		}
		if err = call.Store(handler); err != nil {
			return
		}
	}

	if handler == nil {
		err = ErrDbusOpfailNoHandle
		return
	} else {
		w.handle = *handler
	}

	w.hasHandle = true
	w.IsUnlocked = true

	return
}

/*
	RemoveFolder removes a Folder folderName from a Wallet.
	Note that this will also remove all WalletItems in the given Folder.
*/
func (w *Wallet) RemoveFolder(folderName string) (err error) {

	var call *dbus.Call
	var success bool

	if err = w.walletCheck(); err != nil {
		return
	}

	if call = w.Dbus.Call(
		DbusWMRemoveFolder, 0, w.handle, folderName, w.wm.AppID,
	); call.Err != nil {
		err = call.Err
		return
	}
	if err = call.Store(&success); err != nil {
		return
	}

	if !success {
		err = ErrDbusOpfailRemoveFolder
		return
	}

	return
}

// Update fetches/updates all Folder objects in a Wallet.
func (w *Wallet) Update() (err error) {

	var folderNames []string
	var errs []error = make([]error, 0)

	if err = w.walletCheck(); err != nil {
		return
	}

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
		err = ErrInitWallet
		return
	}

	if _, err = w.IsOpen(); err != nil {
		return
	}

	if !w.IsUnlocked || !w.hasHandle {
		if err = w.Open(); err != nil {
			return
		}
	}

	return
}
