package main

import (
	"fmt"
	"log"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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
	coldResetCard(target)

	_, err = initEmulation(emulator)
	check(err)

	for {
		command, err := receiveCommand(emulator)

		// If we didn't actually receive an APDU, keep waiting.
		if err != nil || len(command) < 3 || command[2] != 0x00 {
			continue
		}

		// Ok, the reader sent `proxiedCommand` to us.
		proxiedCommand := command[3 : len(command)-2]
		log.Println("Received", asHex(proxiedCommand))

		// Let's wrap this DESFire-native command with 7816. This has to happen
		// since something is randomly putting these cards in this mode, and
		// once a 7816 frame is sent, you cannot go back...
		log.Println("Sending the target wrapped; unwrapped:", asHex(proxiedCommand))
		targetResponse, err := normalTransmit(target, wrapCommand(proxiedCommand))
		log.Println("Target responded", len(targetResponse), asHex(targetResponse), err)

		// Move the APDU SW2 to the first byte.
		fixedResponse := targetResponse[:len(targetResponse)-2]
		fixedSw2 := targetResponse[len(targetResponse)-1]
		fixedResponse = append([]byte{fixedSw2}, fixedResponse...)

		// Apple Pay uses DESFire GET VERSION but does not close it out, resulting in
		// COMMAND_ABORTED if not treated. Probably because of our re-framing.
		if proxiedCommand[0] == 0x60 && fixedSw2 == 0xaf {
			normalTransmit(target, wrapCommand([]byte{0xaf}))
			normalTransmit(target, wrapCommand([]byte{0xaf}))
		}

		// Send our fixed response back.
		log.Println("Sending", asHex(fixedResponse), "back.")
		sendResponse(emulator, fixedResponse)
	}
}
