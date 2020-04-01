package util

import (
	"encoding/hex"
	"reflect"
)

// GetAddressFromScriptHash returns base58 encoded address.
func GetAddressFromScriptHash(scriptHash []byte) string {
	if len(scriptHash) == 0 {
		return ""
	}
	scriptHashDecode := append([]byte{0x17}, scriptHash...)
	scriptHashDecode = append(scriptHashDecode, Hash256(scriptHashDecode)[0:4]...)
	addr := EncodeBase58(scriptHashDecode)
	return addr
}

// GetScriptHashFromAddress returns script hash of address.
func GetScriptHashFromAddress(addr string) []byte {
	if len(addr) == 0 {
		return nil
	}
	v, err := DecodeBase58(addr)
	if err != nil {
		panic(err)
	}
	v = ReverseBytes(v)
	addrsc := hex.EncodeToString(v)
	addrsc = addrsc[8 : len(addrsc)-2]
	result, err := hex.DecodeString(addrsc)
	result = ReverseBytes(result)
	if err != nil {
		panic(err)
	}
	return result
}

// AddressValid checks if address is valid.
func AddressValid(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	buffer, err := DecodeBase58(addr)
	if err != nil {
		return false
	}

	if len(buffer) < 4 {
		return false
	}

	checksum := Sha256(Sha256(buffer[:len(buffer)-4]))
	return reflect.DeepEqual(buffer[len(buffer)-4:], checksum[:4])
}

// AddrScValid checks if address script is valid.
func AddrScValid(addrSc string) bool {
	addrScLen := len(addrSc)
	if addrScLen == 0 {
		return true
	}

	if addrScLen != 40 {
		return false
	}

	addrScBytes, _ := hex.DecodeString(addrSc)
	addr := GetAddressFromScriptHash(addrScBytes)

	return AddressValid(addr)
}
