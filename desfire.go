package main

import (
	"github.com/ebfe/scard"
)

// Transforms a native DESFire command like `0x60` into
// a wrapped APDU like `0x90 0x60 0x00 0x00`.
func wrapCommand(command []byte) []byte {
	wrapper := []byte{0x90, command[0], 0x00, 0x00} // CLA, INS, P1, P2 bytes
	if len(command) > 1 {
		wrapper = append(wrapper, uint8(len(command)-1)) // Data length
		wrapper = append(wrapper, command[1:]...)        // Data
	}
	wrapper = append(wrapper, 0x00) // Le byte

	return wrapper
}

// When a DESFire card receives a 7816-4 wrapped frame at any point, it refuses
// to recognize any future DESFire-native commands. As a result we need to
// cold reset the target card until it cooperates.
func ensureNativeCommands(card *scard.Card) {
	for i := 0; i < 3; i++ {
		testRes, err := normalTransmit(card, []byte{0x6a})
		check(err)

		if testRes[0] == 0x67 || testRes[0] == 0x68 {
			// Try cold reset.
			coldResetCard(card)
		} else {
			return
		}
	}

	panic("DESFire card not correctly allowing native command")
}
