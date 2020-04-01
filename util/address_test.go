package util

import (
	"testing"
)

func TestVerifyAddress(t *testing.T) {
	if !AddressValid("APyEx5f4Zm4oCHwFWiSTaph1fPBxZacYVR") {
		t.Error("Address is valid but function [AddressValid] returns invalid result")
	}

	if AddressValid("APyEx5f4Zm4oCHwFWiSTaph1fPBxZacYVA") {
		t.Error("Address is invalid but function [AddressValid] returns valid result")
	}
}
func TestAddrConvertion(t *testing.T) {
	addr := "AKQjaQ7Hor11BfRnXUBvYYiY1CwUkLywyc"
	addrSc := GetScriptHashFromAddress(addr)

	if addr != GetAddressFromScriptHash(addrSc) {
		t.Error("Address convertion failed")
	}
}
