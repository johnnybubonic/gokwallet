// See LICENSE in source root directory for copyright and licensing information.

/*
Package gokwallet serves as a Golang interface to KDE's KWallet (https://utils.kde.org/projects/kwalletmanager/).

Note that this library interfaces with kwalletd. KWallet is in the process of moving to libsecret/SecretService
(see https://bugs.kde.org/show_bug.cgi?id=313216 and https://invent.kde.org/frameworks/kwallet/-/merge_requests/11).
While there is a pull request in place, it has not yet been merged in (and it may be a while before downstream
distributions incorporate that version). However, when that time comes I highly recommend using my `gosecret`
library to interface with that (module r00t2.io/gosecret; see https://pkg.go.dev/r00t2.io/gosecret).

KWallet has the following structure (modified slightly to reflect this library):

- A main Dbus service interface ("org.kde.kwalletd5") allows one to retrieve and operate on/with Wallet items.

- One or more Wallet items allow one to retrieve and operate on/with Folder items.

- One or more Folder items allow one to retrieve and operate on/with Passwords, Maps, BinaryData, and Unknown WalletItem items.

Thus, the hierarchy (as exposed by this library) looks like this:

	WalletManager
	├─ Wallet "A"
	│	├─ Folder "A_1"
	│	│	├─ Passwords
	│	│	├─ Maps
	│	│	├─ BinaryData
	│	│	└─ Unknown
	│	└─ Folder "A_2"
	│		├─ Passwords
	│		├─ Maps
	│		├─ BinaryData
	│		└─ Unknown
	└─ Wallet "B"
		├─ Folder "B_1"
		│	├─ Passwords
		│	├─ Maps
		│	├─ BinaryData
		│	└─ Unknown
		└─ Folder "B_2"
			├─ Passwords
			├─ Maps
			├─ BinaryData
			└─ Unknown

*/
package gokwallet
