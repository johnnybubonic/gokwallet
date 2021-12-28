package gokwallet

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/godbus/dbus/v5"
)

/*
	resultCheck checks the result code from a Dbus call and returns an error if not successful.
	See also resultPassed.
*/
func resultCheck(result int32) (err error) {

	// This is technically way more complex than it needs to be, but is extendable for future use.
	switch i := result; i {
	case DbusSuccess:
		err = nil
	case DbusFailure:
		err = ErrOperationFailed
	default:
		err = ErrOperationFailed
	}

	return
}

/*
	resultPassed checks the result code from a Dbus call and returns a boolean as to whether the result is pass or not.
	See also resultCheck.
*/
func resultPassed(result int32) (passed bool) {

	// This is technically way more complex than it needs to be, but is extendable for future use.
	switch i := result; i {
	case DbusSuccess:
		passed = true
	case DbusFailure:
		passed = false
	default:
		passed = false
	}

	return
}

// bytemapKeys is used to parse out Map names when fetching from Dbus.
func bytemapKeys(variant dbus.Variant) (keyNames []string) {

	var d map[string]dbus.Variant

	d = variant.Value().(map[string]dbus.Variant)

	keyNames = make([]string, len(d))

	idx := 0
	for k, _ := range d {
		keyNames[idx] = k
		idx++
	}

	return
}

// bytesToMap takes a byte slice and returns a map[string]string based on a Dbus QMap struct(ure).
func bytesToMap(raw []byte) (m map[string]string, numEntries uint32, err error) {

	var buf *bytes.Reader
	var kLen uint32
	var vLen uint32
	var k []byte
	var v []byte

	/*
		I considered using:
		- https://github.com/lunixbochs/struc
		- https://github.com/roman-kachanovsky/go-binary-pack
		- https://github.com/go-restruct/restruct

		The second hasn't been updated in quite some time, the first or third would have been a headache due to the variable length,
		and ultimately I felt it was silly to add a dependency for only a single piece of data (Map).
		So sticking to stdlib.
	*/

	buf = bytes.NewReader(raw)

	if err = binary.Read(buf, binary.BigEndian, &numEntries); err != nil {
		return
	}

	m = make(map[string]string, numEntries)

	for i := uint32(0); i < numEntries; i++ {
		if err = binary.Read(buf, binary.BigEndian, &kLen); err != nil {
			return
		}

		k = make([]byte, kLen)

		if err = binary.Read(buf, binary.BigEndian, &k); err != nil {
			return
		}

		if err = binary.Read(buf, binary.BigEndian, &vLen); err != nil {
			return
		}

		v = make([]byte, vLen)

		if err = binary.Read(buf, binary.BigEndian, &v); err != nil {
			return
		}

		// QMap does this infuriating thing where it separates each character with a null byte. So we need to strip them out.
		k = bytes.ReplaceAll(k, []byte{0x0}, []byte{})
		v = bytes.ReplaceAll(v, []byte{0x0}, []byte{})

		m[string(k)] = string(v)
	}

	return
}

// mapToBytes performs the inverse of bytesToMap.
func mapToBytes(m map[string]string) (raw []byte, err error) {

	var numEntries uint32
	var buf *bytes.Buffer
	var kLen uint32
	var vLen uint32
	var kB []byte
	var vB []byte

	if m == nil {
		err = ErrInvalidMap
		return
	}

	numEntries = uint32(len(m))

	buf = &bytes.Buffer{}

	if err = binary.Write(buf, binary.BigEndian, &numEntries); err != nil {
		return
	}

	for k, v := range m {
		kB = []byte(strings.Join(strings.Split(k, ""), "\x00"))
		vB = []byte(strings.Join(strings.Split(v, ""), "\x00"))
		kLen = uint32(len(kB))
		vLen = uint32(len(vB))

		if err = binary.Write(buf, binary.BigEndian, &kLen); err != nil {
			return
		}
		if err = binary.Write(buf, binary.BigEndian, &kB); err != nil {
			return
		}

		if err = binary.Write(buf, binary.BigEndian, &vLen); err != nil {
			return
		}
		if err = binary.Write(buf, binary.BigEndian, &vB); err != nil {
			return
		}
	}

	raw = buf.Bytes()

	return
}
