package gokwallet

import (
	"testing"
)

// TestWallet tests all functions of a Wallet.
func TestWallet(t *testing.T) {

	var err error
	var b bool
	var conns []string
	var folders []string
	var wallets []string
	var r *RecurseOpts = DefaultRecurseOpts
	var wm *WalletManager
	var wm2 *WalletManager
	var w *Wallet
	var w2 *Wallet

	r.AllWalletItems = true

	if wm, err = NewWalletManager(r, appIdTest); err != nil {
		t.Fatalf("failed to get WalletManager '%v' for TestWallet: %v", appIdTest, err)
	}
	defer wm.Close()

	if wm2, err = NewWalletManager(r, appIdTestAlt); err != nil {
		t.Fatalf("failed to get WalletManager '%v' for TestWallet: %v", appIdTest, err)
	}
	defer wm2.Close()

	if w, err = NewWallet(wm, walletTest.String(), r); err != nil {
		t.Fatalf("failed to get Wallet '%v:%v' for TestWallet: %v", appIdTest, walletTest.String(), err)
	}
	defer w.Delete()
	defer w.Disconnect()

	// We test Disconnect above but we also need to test explicit disconnect by application name.
	if w2, err = NewWallet(wm2, walletTestAlt.String(), r); err != nil {
		t.Fatalf("failed to get Wallet '%v:%v' for TestWallet: %v", appIdTestAlt, walletTestAlt.String(), err)
	}
	defer w2.Delete()

	if wallets, err = wm.WalletNames(); err != nil {
		t.Errorf("failure when getting wallet names for '%v': %v", appIdTest, err)
	} else {
		t.Logf("wallet names found via %v: %#v", appIdTest, wallets)
	}

	if w2, err = NewWallet(wm, walletTestAlt.String(), r); err != nil {
		t.Errorf("could not open '%v' in '%v': %v", walletTestAlt.String(), appIdTest, err)
	}
	if err = w2.DisconnectApplication(appIdTestAlt); err != nil {
		t.Errorf(
			"failed to execute DisconnectApplication for '%v:%v' successfully: %v", appIdTestAlt, walletTestAlt.String(), err,
		)
	}

	if err = w.ChangePassword(); err != nil {
		t.Errorf("failed to change password for wallet '%v:%v': %v", appIdTest, walletTest.String(), err)
	}

	if conns, err = w.Connections(); err != nil {
		t.Errorf("failed to get Connections for '%v:%v': %v", appIdTest, walletTest.String(), err)
	} else {
		if conns == nil || len(conns) == 0 {
			t.Errorf("failed to get at least one connection for '%v:%v'. Connections: %#v", appIdTest, walletTest.String(), conns)
		} else {
			t.Logf("Connections for '%v:%v': %#v", appIdTest, walletTest.String(), conns)
		}
	}

	if err = w.CreateFolder(folderTest.String()); err != nil {
		t.Errorf("error when creating folder '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err)
	} else {
		t.Logf("created folder '%v:%v:%v'", appIdTest, walletTest.String(), folderTest.String())
	}

	if b, err = w.FolderExists(folderTest.String()); err != nil {
		t.Errorf(
			"error when running FolderExists for '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err,
		)
	} else if !b {
		t.Errorf(
			"did not detecting existing folder '%v:%v:%v' in FolderExists", appIdTest, walletTest.String(), folderTest.String(),
		)
	}
	if b, err = w.HasFolder(folderTest.String()); err != nil {
		t.Errorf(
			"error when running HasFolder for '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err,
		)
	} else if !b {
		t.Errorf(
			"did not detecting existing folder '%v:%v:%v' in HasFolder", appIdTest, walletTest.String(), folderTest.String(),
		)
	}

	if folders, err = w.ListFolders(); err != nil {
		t.Errorf("error when running ListFolders for wallet '%v:%v': %v", appIdTest, walletTest.String(), err)
	} else {
		t.Logf("ListFolders returned for wallet '%v:%v': %#v", appIdTest, walletTest.String(), folders)
	}

	if err = w.RemoveFolder(folderTest.String()); err != nil {
		t.Errorf("failed running RemoveFolder in Wallet for '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err)
	}
}
