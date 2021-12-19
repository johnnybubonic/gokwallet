// See LICENSE in source root directory for copyright and licensing information.

/*
Package gokwallet serves as a Golang interface to KDE's KWallet (https://utils.kde.org/projects/kwalletmanager/).

Note that to use this library, the running machine must have both Dbus and kwalletd running.

Relatedly, note also that this library interfaces with kwalletd. KWallet is in the process of moving to libsecret/SecretService
(see https://bugs.kde.org/show_bug.cgi?id=313216 and https://invent.kde.org/frameworks/kwallet/-/merge_requests/11),
thus replacing kwalletd.
While there is a pull request in place, it has not yet been merged in (and it may be a while before downstream
distributions incorporate that version). However, when that time comes I highly recommend using my `gosecret`
library to interface with that (module r00t2.io/gosecret; see https://pkg.go.dev/r00t2.io/gosecret).

KWallet has the following structure (modified slightly to reflect this library):

- A main Dbus service interface ("org.kde.kwalletd5"), WalletManager, allows one to retrieve and operate on/with Wallet items.

- One or more Wallet items allow one to retrieve and operate on/with Folder items.

- One or more Folder items allow one to retrieve and operate on/with Passwords, Maps, BinaryData, and Unknown WalletItem items.

Thus, the hierarchy (as exposed by this library) looks like this:

	WalletManager
	├─ Wallet "A"
	│	├─ Folder "A_1"
	│	│	├─ Passwords
	│	│	│	├─ Password "A_1_a"
	│	│	│	└─ Password "A_1_b"
	│	│	├─ Maps
	│	│	│	├─ Map "A_1_a"
	│	│	│	└─ Map "A_1_b"
	│	│	├─ BinaryData
	│	│	│	├─ Blob "A_1_a"
	│	│	│	└─ Blob "A_1_b"
	│	│	└─ Unknown
	│	│		├─ UnknownItem "A_1_a"
	│	│		└─ UnknownItem "A_1_b"
	│	└─ Folder "A_2"
	│		├─ Passwords
	│		│	├─ Password "A_2_a"
	│		│	└─ Password "A_2_b"
	│		├─ Maps
	│		│	├─ Map "A_2_a"
	│		│	└─ Map "A_2_b"
	│		├─ BinaryData
	│		│	├─ Blob "A_2_a"
	│		│	└─ Blob "A_2_b"
	│		└─ Unknown
	│			├─ UnknownItem "A_2_a"
	│			└─ UnknownItem "A_2_b"
	└─ Wallet "B"
		└─ Folder "B_1"
			├─ Passwords
			│	├─ Password "B_1_a"
			│	└─ Password "B_1_b"
			├─ Maps
			│	├─ Map "B_1_a"
			│	└─ Map "B_1_b"
			├─ BinaryData
			│	├─ Blob "B_1_a"
			│	└─ Blob "B_1_b"
			└─ Unknown
				├─ UnknownItem "B_1_a"
				└─ UnknownItem "B_1_b"

This is an approximation, but should show a relatively accurate representation of the model.
Note that most systems are likely to only have a single wallet, "kdewallet".

Usage

Full documentation can be found via inline documentation.
Additionally, use either https://pkg.go.dev/r00t2.io/gokwallet or https://pkg.go.dev/golang.org/x/tools/cmd/godoc (or `go doc`) in the source root.

You most likely do *not* want to call any New<object> function directly;
NewWalletManager with its RecurseOpts parameter (`recursion`) should get you everything you want/need.
*/
package gokwallet
