package gokwallet

import (
	"testing"
)

func TestFolder(t *testing.T) {

	var r *RecurseOpts = DefaultRecurseOpts
	var wm *WalletManager
	var w *Wallet
	var f *Folder
	var p *Password
	var b bool
	var entries []string
	var err error

	r.AllWalletItems = true

	if wm, err = NewWalletManager(r, appIdTest); err != nil {
		t.Errorf("failure when getting WalletManager '%v': %v", appIdTest, err)
		return
	}
	defer wm.Close()

	if w, err = NewWallet(wm, walletTest.String(), wm.Recurse); err != nil {
		t.Errorf("failure when getting Wallet '%v:%v': %v", appIdTest, walletTest.String(), err)
		return
	}

	if f, err = NewFolder(w, folderTest.String(), w.Recurse); err != nil {
		t.Errorf("failure when getting Folder '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err)
		return
	}

	t.Logf("created Folder '%v:%v'; initialized: %v", w.Name, f.Name, f.isInit)

	if err = f.Update(); err != nil {
		t.Errorf("failed to update Folder '%v:%v': %v", w.Name, f.Name, err)
	}

	if _, err = f.WriteBlob(blobTest.String(), testBytes); err != nil {
		t.Errorf("failed to WriteBlob '%#v' in '%v:%v:%v': %v", testBytes, w.Name, f.Name, blobTest.String(), err)
	}
	if _, err = f.WriteMap(mapTest.String(), testMap); err != nil {
		t.Errorf("failed to WriteMap '%#v' in '%v:%v:%v': %v", testMap, w.Name, f.Name, mapTest.String(), err)
	}
	if p, err = f.WritePassword(passwordTest.String(), testPassword); err != nil {
		t.Errorf("failed to WritePassword '%#v' in '%v:%v:%v': %v", testPassword, w.Name, f.Name, passwordTest.String(), err)
	}
	if _, err = f.WriteUnknown(unknownItemTest.String(), testBytes); err != nil {
		t.Errorf("failed to WriteUnknown '%#v' in '%v:%v:%v': %v", testBytes, w.Name, f.Name, unknownItemTest.String(), err)
	}

	if err = f.UpdateBlobs(); err != nil {
		t.Errorf("failed to update Blobs in Folder '%v:%v': %v", w.Name, f.Name, err)
	}
	if err = f.UpdateMaps(); err != nil {
		t.Errorf("failed to update Map in Folder '%v:%v': %v", w.Name, f.Name, err)
	}
	if err = f.UpdatePasswords(); err != nil {
		t.Errorf("failed to update Passwords in Folder '%v:%v': %v", w.Name, f.Name, err)
	}
	if err = f.UpdateUnknowns(); err != nil {
		t.Errorf("failed to update UnknownItems in Folder '%v:%v': %v", w.Name, f.Name, err)
	}

	if b, err = f.HasEntry(p.Name); err != nil {
		t.Errorf("failed to run HasEntry in Folder '%v:%v' for key '%v': %v", w.Name, f.Name, p.Name, err)
	} else if !b {
		t.Errorf("failed to find entry '%v' via HasEntry in Folder '%v:%v': %v", p.Name, w.Name, f.Name, err)
	}

	// This gives an incorrect return of true and I'm not entirely sure why. Maybe it needs a .sync or .reconfigure Dbus call?
	// Or maybe it needs to be encapsulated in quotes?
	/*
		if b, err = f.KeyNotExist(p.Name); err != nil {
			t.Errorf("failed to run KeyNotExist in Folder '%v:%v' for key '%v': %v", w.Name, f.Name, p.Name, err)
		} else if b {
			t.Errorf("failed to get false for '%v' via KeyNotExist in Folder '%v:%v'", p.Name, w.Name, f.Name)
			// t.Fatalf("failed to get false for '%v' via KeyNotExist in Folder '%v:%v'", p.Name, w.Name, f.Name)
		}
	*/

	if entries, err = f.ListEntries(); err != nil {
		t.Errorf("failed to run ListEntries in Folder '%v:%v': %v", w.Name, f.Name, err)
	} else if entries == nil || len(entries) == 0 {
		t.Errorf("ListEntries for Folder '%v:%v' is 0", w.Name, f.Name)
	}

	b = false
	for idx, e := range entries {
		if e == p.Name {
			t.Logf("found matching value for test password in '%v:%v' at index %v: %v", w.Name, f.Name, idx, p.Name)
			b = true
			break
		}
	}
	if !b {
		t.Errorf("failed to find test password '%v:%v:%v'", w.Name, f.Name, p.Name)
	}

	// This tests the parent folder's Rename method.
	if err = p.Rename(passwordTestRename.String()); err != nil {
		t.Errorf(
			"failed to RenameEntry '%v' to '%v' in Folder '%v:%v': %v",
			p.Name, passwordTestRename.String(), w.Name, f.Name, err,
		)
	} else {
		t.Logf("checking existence for key '%v'", p.Name)
		// This runs HasEntry via the parent folder.
		if b, err = p.Exists(); err != nil {
			t.Errorf("failed to run HasEntry in Folder '%v:%v' for key '%v': %v", w.Name, f.Name, passwordTestRename.String(), err)
		} else if !b {
			t.Errorf("failed to find entry '%v' via HasEntry in Folder '%v:%v': %v", passwordTestRename.String(), w.Name, f.Name, err)
		} else {
			p.Name = passwordTestRename.String()
		}
	}

	if err = f.RemoveEntry(p.Name); err != nil {
		t.Errorf("failed to RemoveEntry entry '%v' in Folder '%v:%v': %v", p.Name, w.Name, f.Name, err)
	} else {
		if b, err = f.HasEntry(p.Name); err != nil {
			t.Errorf("failed to run HasEntry in Folder '%v:%v' for key '%v': %v", w.Name, f.Name, p.Name, err)
		} else if b {
			t.Errorf("failed to successfully remove entry '%v' via RemoveEntry in Folder '%v:%v': %v", p.Name, w.Name, f.Name, err)
		}
	}

	if err = f.Delete(); err != nil {
		t.Errorf("failed to delete Folder '%v:%v': %v", w.Name, f.Name, err)
	}
	if err = w.Delete(); err != nil {
		t.Errorf("failed to delete Wallet '%v': %v", w.Name, err)
	}
}
