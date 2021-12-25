package gokwallet

import (
	"testing"
)

func getTestEnv(t *testing.T) (wm *WalletManager, w *Wallet, f *Folder, err error) {

	var r *RecurseOpts = DefaultRecurseOpts

	r.AllWalletItems = true

	if wm, err = NewWalletManager(r, appIdTest); err != nil {
		t.Errorf("failure when getting WalletManager '%v': %v", appIdTest, err)
		return
	}

	if w, err = NewWallet(wm, walletTest.String(), wm.Recurse); err != nil {
		t.Errorf("failure when getting Wallet '%v:%v': %v", appIdTest, walletTest.String(), err)
		return
	}

	if f, err = NewFolder(w, folderTest.String(), w.Recurse); err != nil {
		t.Errorf("failure when getting Folder '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err)
		return
	}

	return
}
