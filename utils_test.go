package gokwallet

import (
	"testing"
)

func getTestEnv(t *testing.T) (e *testEnv, err error) {

	e = &testEnv{
		wm: nil,
		w:  nil,
		f:  nil,
		r:  DefaultRecurseOpts,
	}

	e.r.AllWalletItems = true

	if e.wm, err = NewWalletManager(e.r, appIdTest); err != nil {
		t.Errorf("failure when getting WalletManager '%v': %v", appIdTest, err)
		return
	}

	if e.w, err = NewWallet(e.wm, walletTest.String(), e.r); err != nil {
		t.Errorf("failure when getting Wallet '%v:%v': %v", appIdTest, walletTest.String(), err)
		return
	}

	if e.f, err = NewFolder(e.w, folderTest.String(), e.r); err != nil {
		t.Errorf("failure when getting Folder '%v:%v:%v': %v", appIdTest, walletTest.String(), folderTest.String(), err)
		return
	}

	return
}

// cleanup closes connections and deletes created folders/wallets in a testEnv.
func (e *testEnv) cleanup(t *testing.T) (err error) {

	var errs []error = make([]error, 0)

	if err = e.f.Delete(); err != nil {
		errs = append(errs, err)
		err = nil
	}

	if err = e.w.Delete(); err != nil {
		errs = append(errs, err)
		err = nil
	}

	if err = e.wm.Close(); err != nil {
		errs = append(errs, err)
		err = nil
	}

	if errs != nil && len(errs) > 0 {
		err = NewErrors(errs...)
		return
	}

	return
}
