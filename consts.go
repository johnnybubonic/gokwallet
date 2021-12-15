package gokwallet

// KWalletD Dbus interfaces.
const (
	// DbusService is the Dbus service bus identifier.
	DbusService string = "org.kde.kwalletd5"
	// DbusServiceBase is the base identifier used by interfaces.
	DbusServiceBase string = "org.kde"
)

// gokwallet defaults.
const (
	// DefaultWalletName is the name of the default Wallet to use.
	DefaultWalletName string = "kdewallet"
	// DefaultAppID is the default name for the application (see WalletManager.AppID).
	DefaultAppID string = "GoKwallet"
)

// WalletManager interface.
const (
	/*
		DbusInterfaceWM is the Dbus interface for working with a WalletManager.
	*/
	DbusInterfaceWM string = DbusServiceBase + ".KWallet"

	// Methods

	// DbusWMChangePassword changes the password for a Wallet.
	DbusWMChangePassword string = DbusInterfaceWM + ".changePassword"

	// DbusWMClose closes an App (WalletManager) or Wallet.
	DbusWMClose string = DbusInterfaceWM + ".close"

	// DbusWMCloseAllWallets closes all WalletManager.Wallets.
	DbusWMCloseAllWallets string = DbusInterfaceWM + ".closeAllWallets"

	// DbusWMCreateFolder creates a Folder.
	DbusWMCreateFolder string = DbusInterfaceWM + ".createFolder"

	// DbusWMDeleteWallet deletes/removes a Wallet.
	DbusWMDeleteWallet string = DbusInterfaceWM + ".deleteWallet"

	// DbusWMDisconnectApp disconnects a WalletManager (or other App).
	DbusWMDisconnectApp string = DbusInterfaceWM + ".disconnectApplication"

	// DbusWMEntriesList returns a *map* of the WalletItem objects in a Folder (with their entry name as the map key).
	DbusWMEntriesList string = DbusInterfaceWM + ".entriesList"

	// DbusWMEntryList returns a *slice* of WalletItem names in a Folder.
	DbusWMEntryList string = DbusInterfaceWM + ".entryList"

	// DbusWMEntryType returns the type of a WalletItem.
	DbusWMEntryType string = DbusInterfaceWM + ".entryType"

	// DbusWMFolderNotExist indicates if a Folder exists within a Wallet or not.
	DbusWMFolderNotExist string = DbusInterfaceWM + ".folderDoesNotExist"

	// DbusWMFolderList lists the Folder objects (as Folder.Name) in a Wallet.
	DbusWMFolderList string = DbusInterfaceWM + ".folderList"

	// DbusWMHasEntry indicates if a Folder has a WalletItem or not.
	DbusWMHasEntry string = DbusInterfaceWM + ".hasEntry"

	// DbusWMHasFolder indicates if a Wallet has a Folder or not.
	DbusWMHasFolder string = DbusInterfaceWM + ".hasFolder"

	/*
		DbusWMIsEnabled indicates if KWallet is enabled.
		TODO: Is this accurate?
	*/
	DbusWMIsEnabled string = DbusInterfaceWM + ".isEnabled"

	// DbusWMIsOpen indicates if a Wallet is open (unlocked).
	DbusWMIsOpen string = DbusInterfaceWM + ".isOpen"

	// DbusWMKeyNotExist indicates if a Folder has a WalletItem or not.
	DbusWMKeyNotExist string = DbusInterfaceWM + ".keyDoesNotExist"

	// DbusWMLocalWallet gives the name of the local (default?) Wallet.
	DbusWMLocalWallet string = DbusInterfaceWM + ".localWallet"

	// DbusWMMapList gives a list of Map names in a Folder.
	DbusWMMapList string = DbusInterfaceWM + ".mapList"

	/*
		DbusWMNetWallet indicates if a Wallet is a Network Wallet or not.
		TODO: is/was this ever used?
	*/
	DbusWMNetWallet string = DbusInterfaceWM + ".networkWallet"

	// DbusWMOpen opens (unlocks) a Wallet.
	DbusWMOpen string = DbusInterfaceWM + ".open"

	// DbusWMOpenAsync opens (unlocks) a Wallet asynchronously.
	DbusWMOpenAsync string = DbusInterfaceWM + ".openAsync"

	// DbusWMOpenPath opens (unlocks) a Wallet by its filepath.
	DbusWMOpenPath string = DbusInterfaceWM + ".openPath"

	// DbusWMOpenPathAsync opens (unlocks) a Wallet by its filepath asynchronously.
	DbusWMOpenPathAsync string = DbusInterfaceWM + ".openPath"

	// DbusWMPamOpen opens (unlocks) a Wallet via PAM.
	DbusWMPamOpen string = DbusInterfaceWM + ".pamOpen"

	/*
		DbusWMPasswordList returns a map of Password objects in a Folder.
		Password.Name is the map key.
	*/
	DbusWMPasswordList string = DbusInterfaceWM + ".passwordList"

	// DbusWMReadEntry fetches a WalletItem by its name from a Folder (as a byteslice).
	DbusWMReadEntry string = DbusInterfaceWM + ".readEntry"

	// DbusWMReadEntryList returns a map of WalletItem objects in a Folder.
	DbusWMReadEntryList string = DbusInterfaceWM + ".readEntryList"

	// DbusWMReadMap returns a Map from a Folder (as a byteslice).
	DbusWMReadMap string = DbusInterfaceWM + ".readMap"

	// DbusWMReadMapList returns a map of Map objects in a Folder.
	DbusWMReadMapList string = DbusInterfaceWM + ".readMapList"

	// DbusWMReadPassword returns a Password from a Folder (as a byteslice).
	DbusWMReadPassword string = DbusInterfaceWM + ".readPassword"

	// DbusWMReadPasswordList returns a map of Password objects in a Folder.
	DbusWMReadPasswordList string = DbusInterfaceWM + ".readPasswordList"

	// DbusWMReconfigure is [FUNCTION UNKNOWN/UNDOCUMENTED; TODO? NOT IMPLEMENTED.]
	// DbusWMReconfigure string = DbusInterfaceWM + ".reconfigure"

	// DbusWMRemoveEntry removes a WalletItem from a Folder.
	DbusWMRemoveEntry string = DbusInterfaceWM + ".removeEntry"

	// DbusWMRemoveFolder removes a Folder from a Wallet.
	DbusWMRemoveFolder string = DbusInterfaceWM + ".removeFolder"

	// DbusWMRenameEntry renames ("moves") a WalletItem.
	DbusWMRenameEntry string = DbusInterfaceWM + ".renameEntry"

	// DbusWMSync is [FUNCTION UNKNOWN/UNDOCUMENTED; TODO? RELATED TO ASYNC? NOT IMPLEMENTED.]
	// DbusWMSync string = DbusInterfaceWM + ".sync"

	// DbusWMUsers returns a slice of users.
	DbusWMUsers string = DbusInterfaceWM + ".users"

	// DbusWMWallets returns an array of Wallet names.
	DbusWMWallets string = DbusInterfaceWM + ".wallets"

	// DbusWMWriteEntry writes (creates) a WalletItem to/in a Folder.
	DbusWMWriteEntry string = DbusInterfaceWM + ".writeEntry"

	// DbusWMWriteMap writes (creates) a Map (via a byteslice) to/in a Folder.
	DbusWMWriteMap string = DbusInterfaceWM + ".writeMap"

	// DbusWMWritePassword writes (creates) a Password to/in a Folder.
	DbusWMWritePassword string = DbusInterfaceWM + ".writePassword"
)

// Dbus paths.
const (
	// DbusPath is the path for DbusService.
	DbusPath string = "/modules/kwalletd5"
)
