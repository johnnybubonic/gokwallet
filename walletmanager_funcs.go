package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewWalletManager returns a WalletManager. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	If appId is empty/nil, DefaultAppID will be used as the app ID.
	If appId is specified, only the first string is used.
*/
func NewWalletManager(recursion *RecurseOpts, appID ...string) (wm *WalletManager, err error) {

	var realAppID string

	if appID != nil && len(appID) > 0 {
		realAppID = appID[0]
	} else {
		realAppID = DefaultAppID
	}

	if wm, err = newWM(realAppID, recursion); err != nil {
		return
	}

	return
}

/*
	NewWalletManagerFiles returns a WalletManager from one or more filePaths (*.kwl, *.tar, or *.xml exports).
	Note that if the wallet file was created via an "encrypted export", it will be a .kwl file
	inside a .tar.
	err will contain a MultiError if any filepaths specified do not exist or cannot be opened.
	It requires a RecurseOpts (you can use DefaultRecurseOpts, call NewRecurseOpts,
	or provide your own RecurseOpts struct).
	If appId is empty, DefaultAppID will be used as the app ID.
*/
/* TODO: POC this before exposing. I have NO idea if/how it'll work.
func NewWalletManagerFiles(recursion *RecurseOpts, appId string, filePaths ...string) (wm *WalletManager, err error) {

	var exist bool
	var errs []error = make([]error, 0)
	var realFilePaths []string = make([]string, 0)

	if appId == "" {
		appId = DefaultAppID
	}

	for _, f := range filePaths {
		if f == "" {
			continue
		}
		if exist, err = paths.RealPathExists(&f); err != nil {
			errs = append(errs, err)
			err = nil
			continue
		}
		if !exist {
			err = errors.New(fmt.Sprintf("%v does not exist", f))
			err = nil
			continue
		}
		realFilePaths = append(realFilePaths, f)
	}

	// TODO: do the actual newWM here.

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
	}

	return
}
*/

/*
	CloseWallet closes a Wallet.
	Unlike Wallet.Close, this closes access for ALL applications/WalletManagers
	for the specified Wallet - not just this WalletManager.
*/
func (wm *WalletManager) CloseWallet(walletName string) (err error) {
	var rslt int32

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if err = wm.Dbus.Call(
		DbusWMClose, 0, walletName, false,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

/*
	ForceCloseWallet is like WalletManager.CloseWallet but will still close a Wallet even if currently in use.
	Unlike Wallet.ForceClose, this closes access for ALL applications/WalletManagers
	for the specified Wallet - not just this WalletManager.
*/
func (wm *WalletManager) ForceCloseWallet(walletName string) (err error) {

	var rslt int32

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	// Using a handler allows us to close access for this particular parent WalletManager.
	if err = wm.Dbus.Call(
		DbusWMClose, 0, walletName, false,
	).Store(&rslt); err != nil {
		return
	}

	err = resultCheck(rslt)

	return
}

// CloseAllWallets closes all Wallet objects. They do *not* need to be part of WalletManager.Wallets.
func (wm *WalletManager) CloseAllWallets() (err error) {

	var call *dbus.Call

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	call = wm.Dbus.Call(
		DbusWMCloseAllWallets, 0,
	)
	err = call.Err

	return
}

// IsEnabled returns whether KWallet is enabled or not (and also updates WalletManager.Enabled).
func (wm *WalletManager) IsEnabled() (enabled bool, err error) {

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	if err = wm.Dbus.Call(
		DbusWMIsEnabled, 0,
	).Store(&wm.Enabled); err != nil {
		return
	}

	enabled = wm.Enabled

	return
}

// LocalWallet returns the "local" wallet (and updates WalletManager.Local).
func (wm *WalletManager) LocalWallet() (w *Wallet, err error) {

	var wn string

	if err = wm.Dbus.Call(
		DbusWMLocalWallet, 0,
	).Store(&wn); err != nil {
		return
	}

	if w, err = NewWallet(wm, wn, wm.Recurse); err != nil {
		return
	}

	wm.Local = w

	return
}

// NetworkWallet returns the "network" wallet (and updates WalletManager.Network).
func (wm *WalletManager) NetworkWallet() (w *Wallet, err error) {

	var wn string

	if err = wm.Dbus.Call(
		DbusWMNetWallet, 0,
	).Store(&wn); err != nil {
		return
	}

	if w, err = NewWallet(wm, wn, wm.Recurse); err != nil {
		return
	}

	wm.Network = w

	return
}

// WalletNames returns a list of existing Wallet names.
func (wm *WalletManager) WalletNames() (wallets []string, err error) {

	if err = wm.Dbus.Call(
		DbusWMWallets, 0,
	).Store(&wallets); err != nil {
		return
	}

	return
}

// Update fetches/updates all Wallet objects in a WalletManager.
func (wm *WalletManager) Update() (err error) {

	var walletNames []string
	var errs []error = make([]error, 0)

	if !wm.isInit {
		err = ErrNotInitialized
		return
	}

	if walletNames, err = wm.WalletNames(); err != nil {
		return
	}

	wm.Wallets = make(map[string]*Wallet)

	for _, wn := range walletNames {

		if wm.Wallets[wn], err = NewWallet(wm, wn, wm.Recurse); err != nil {
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

// newWM is what does the heavy lifting behind NewWalletManager and NewWalletManagerFiles.
func newWM(appId string, recursion *RecurseOpts, filePaths ...string) (wm *WalletManager, err error) {

	wm = &WalletManager{
		DbusObject: &DbusObject{
			Conn: nil,
			Dbus: nil,
		},
		AppID:   appId,
		Wallets: nil,
		Recurse: recursion,
	}

	if wm.DbusObject.Conn, err = dbus.SessionBus(); err != nil {
		return
	}
	wm.DbusObject.Dbus = wm.DbusObject.Conn.Object(DbusService, dbus.ObjectPath(DbusPath))

	wm.isInit = true

	if wm.Recurse.All || wm.Recurse.Wallets {
		if err = wm.Update(); err != nil {
			return
		}
		if _, err = wm.LocalWallet(); err != nil {
			return
		}
		if _, err = wm.NetworkWallet(); err != nil {
			return
		}
	}

	wm.isInit = true

	return
}
