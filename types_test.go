package gokwallet

/*
	testEnv is an environment to use for tests.
	It's returned by getTestEnv.
*/
type testEnv struct {
	wm *WalletManager
	w  *Wallet
	f  *Folder
	r  *RecurseOpts
}
