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
		if err != nil || len(command) < 3 {
			continue
		}

		if command[2] != 0x00 {
			continue
		}

		proxiedCommand := command[3 : len(command)-2]
		log.Println("Received", asHex(proxiedCommand))

		log.Println("Sending the target unwrapped:", asHex(proxiedCommand))
		targetResponse, err := normalTransmit(target, wrapCommand(proxiedCommand))
		log.Println("Target responded", len(targetResponse), asHex(targetResponse), err)

		fixedResponse := targetResponse[:len(targetResponse)-2]
		fixedSw2 := targetResponse[len(targetResponse)-1]
		fixedResponse = append([]byte{fixedSw2}, fixedResponse...)

		// Apple Pay uses DESFire GET VERSION but does not close it out, resulting in
		// COMMAND_ABORTED if not treated. Probably because of our re-framing.
		if proxiedCommand[0] == 0x60 && fixedSw2 == 0xaf {
			normalTransmit(target, wrapCommand([]byte{0xaf}))
			normalTransmit(target, wrapCommand([]byte{0xaf}))
		}

		log.Println("Sending", asHex(fixedResponse), "back.")

		sendResponse(emulator, fixedResponse)
	}
}
