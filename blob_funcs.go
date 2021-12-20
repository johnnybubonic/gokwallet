package gokwallet

/*
	NewBlob returns a Blob. It requires a RecurseOpts
	(you can use DefaultRecurseOpts, call NewRecurseOpts, or provide your own RecurseOpts struct).
	It also requires a Folder.
*/
func NewBlob(f *Folder, keyName string, recursion *RecurseOpts) (blob *Blob, err error) {

	if !f.isInit {
		err = ErrNotInitialized
		return
	}

	blob = &Blob{
		DbusObject: f.DbusObject,
		Name:       keyName,
		// Value:      "",
		Recurse: recursion,
		wm:      f.wallet.wm,
		wallet:  f.wallet,
		folder:  f,
		isInit:  false,
	}

	if blob.Recurse.AllWalletItems || blob.Recurse.Blobs {
		if err = blob.Update(); err != nil {
			return
		}
	}

	blob.isInit = true

	return
}

// Update fetches a Blob's Blob.Value.
func (b *Blob) Update() (err error) {

	// TODO.

	return
}

// isWalletItem is needed for interface membership.
func (b *Blob) isWalletItem() (isWalletItem bool) {

	isWalletItem = true

	return
}
