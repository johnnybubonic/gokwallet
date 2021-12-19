package gokwallet

import (
	"errors"
)

/*
	NewRecurseOpts returns a RecurseOpts based on the specified options.
	See the documentation for RecurseOpts for descriptions of the behaviour for each recursion option.
	warn is a MultiError but should be treated as warnings rather than strictly errors.
*/
func NewRecurseOpts(recurseAll, wallets, folders, recurseAllWalletItems, passwords, maps, blobs, unknownItems bool) (opts *RecurseOpts, warn error) {

	var err []error = make([]error, 0)

	if recurseAll {
		if !wallets {
			err = append(err, errors.New("wallets was specified as false but recurseAll is true; recurseAll overrides wallets to true"))
			wallets = true
		}
		if !folders {
			err = append(err, errors.New("folders was specified as false but recurseAll is true; recurseAll overrides folders to true"))
			folders = true
		}
		// NOTE: This is commented out because we control these explicitly with recurseAllWalletItems.
		/*
			if !passwords {
				err = append(err, errors.New("passwords was specified as false but recurseAll is true; recurseAll overrides passwords to true"))
				passwords = true
			}
			if !maps {
				err = append(err, errors.New("maps was specified as false but recurseAll is true; recurseAll overrides maps to true"))
				maps = true
			}
			if !blobs {
				err = append(err, errors.New("blobs was specified as false but recurseAll is true; recurseAll overrides blobs to true"))
				blobs = true
			}
			if !unknownItems {
				err = append(err, errors.New("unknownItems was specified as false but recurseAll is true; recurseAll overrides unknownItems to true"))
				unknownItems = true
			}
			if !recurseAllWalletItems {
				err = append(
					err,
					errors.New(
						"recurseAllWalletItems was specified as false but recurseAll is true; recurseAll overrides recurseAllWalletItems to true",
					),
				)
				recurseAllWalletItems = true
			}
		*/
	} else {
		if recurseAllWalletItems {
			if !passwords {
				err = append(
					err,
					errors.New(
						"passwords was specified as false but recurseAllWalletItems is true; recurseAllWalletItems overrides passwords to true",
					),
				)
				passwords = true
			}
			if !maps {
				err = append(
					err,
					errors.New(
						"maps was specified as false but recurseAllWalletItems is true; recurseAllWalletItems overrides maps to true",
					),
				)
				maps = true
			}
			if !blobs {
				err = append(
					err,
					errors.New(
						"blobs was specified as false but recurseAllWalletItems is true; recurseAllWalletItems overrides blobs to true",
					),
				)
				blobs = true
			}
			if !unknownItems {
				err = append(
					err,
					errors.New(
						"unknownItems was specified as false but recurseAllWalletItems is true; recurseAllWalletItems overrides unknownItems to true",
					),
				)
				unknownItems = true
			}
		}
	}

	opts = &RecurseOpts{
		All:            recurseAll,
		Wallets:        wallets,
		Folders:        folders,
		AllWalletItems: recurseAllWalletItems,
		Passwords:      passwords,
		Maps:           maps,
		Blobs:          blobs,
		UnknownItems:   unknownItems,
	}

	if err != nil && len(err) > 0 {
		warn = NewErrors(err...)
	}

	return
}
