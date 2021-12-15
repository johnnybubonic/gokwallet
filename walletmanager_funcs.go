package gokwallet

import (
	"github.com/godbus/dbus/v5"
)

/*
	NewWalletManager returns a WalletManager.
	If appId is empty/nil, DefaultAppID will be used as the app ID.
	If appId is specified, only the first string is used.
*/
func NewWalletManager(appID ...string) (wm *WalletManager, err error) {

	var realAppID string

	if appID != nil && len(appID) > 0 {
		realAppID = appID[0]
	} else {
		realAppID = DefaultAppID
	}

	wm = &WalletManager{
		DbusObject: &DbusObject{
			Conn: nil,
			Dbus: nil,
		},
		AppID:   realAppID,
		Wallets: make(map[string]*Wallet),
	}

	if wm.DbusObject.Conn, err = dbus.SessionBus(); err != nil {
		return
	}
	wm.DbusObject.Dbus = wm.DbusObject.Conn.Object(DbusService, dbus.ObjectPath(DbusPath))

	return
}

/*
	Update fetches/updates all Wallet objects in a WalletManager.
*/
func (wm *WalletManager) Update() (err error) {

	var wallets []*Wallet

	// TODO.
	_ = wallets

	return
}
