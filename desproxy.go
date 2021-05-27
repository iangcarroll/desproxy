package main

import (
	"fmt"
	"log"
)

// Calls `panic` when an error is present.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns a []byte as a hex string; []byte{0xff, 0xaa} = "ffaa".
func asHex(in []byte) (out string) {
	for _, byt := range in {
		out += fmt.Sprintf("%02x", byt)
	}

	return out
}

func main() {
	// Establish the emulator.
	emulator, err := connectToCard(0, true)
	check(err)

	// Establish the target.
	target, err := connectToCard(1, false)
	check(err)

	// Ensure the target will support native DESFire APDUs.
	ensureNativeCommands(target)

	// Initialize the emulation with the ACR122u.
	_, err = initEmulation(emulator)
	check(err)

	for {
		// Receive a command sent from a reader.
		command, err := receiveCommand(emulator)

		// We lost connection; bring the emulation back up.
		if len(command) > 2 && (command[2] == 0x13 || command[2] == 0x25) {
			log.Println("New emulation session.")

			// Re-initialize the emulation on the emulator.
			_, err = initEmulation(emulator)
			check(err)

			// Cold reset the target card.
			coldResetCard(target)

			// Go back to trying to receive a message.
			continue
		}

		// If we didn't actually receive an APDU, keep waiting.
		if err != nil || len(command) < 3 || command[2] != 0x00 {
			continue
		}

		// Ok, the reader sent `proxiedCommand` to us.
		proxiedCommand := command[3 : len(command)-2]
		log.Println("Received", asHex(proxiedCommand))

		// Send `proxiedCommand` to the target and get the `targetResponse`.
		log.Println("Sending the target:", asHex(proxiedCommand))
		targetResponse, err := normalTransmit(target, proxiedCommand)
		log.Println("Target responded", len(targetResponse), asHex(targetResponse), err)

		// This is the ideal place to manipulate any response data before it goes back to the real reader.

		// Send our fixed response back.
		log.Println("Sending", asHex(targetResponse), "back.")
		sendResponse(emulator, targetResponse)
	}
}
