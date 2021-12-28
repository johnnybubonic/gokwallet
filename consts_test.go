package gokwallet

import (
	"github.com/google/uuid"
)

// Strings.
const (
	appIdTest    string = "GoKwallet_Test"
	appIdTestAlt string = "GoKwallet_Test_Alternate"
)

// Identifiers/names/keys.
var (
	walletTest         uuid.UUID = uuid.New()
	walletTestAlt      uuid.UUID = uuid.New()
	folderTest         uuid.UUID = uuid.New()
	blobTest           uuid.UUID = uuid.New()
	mapTest            uuid.UUID = uuid.New()
	passwordTest       uuid.UUID = uuid.New()
	passwordTestRename uuid.UUID = uuid.New()
	unknownItemTest    uuid.UUID = uuid.New()
)

// Values.
var (
	testBytes           []byte            = []byte(uuid.New().String())
	testBytesReplace    []byte            = []byte(uuid.New().String())
	testPassword        string            = uuid.New().String()
	testPasswordReplace string            = uuid.New().String()
	testMap             map[string]string = map[string]string{
		uuid.New().String(): uuid.New().String(),
		uuid.New().String(): uuid.New().String(),
		uuid.New().String(): uuid.New().String(),
	}
	testMapReplace map[string]string = map[string]string{
		uuid.New().String(): uuid.New().String(),
		uuid.New().String(): uuid.New().String(),
		uuid.New().String(): uuid.New().String(),
	}
)
