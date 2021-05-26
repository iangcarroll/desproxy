package main

func wrapCommand(command []byte) []byte {
	wrapper := []byte{0x90, command[0], 0x00, 0x00} // CLA, INS, P1, P2 bytes
	if len(command) > 1 {
		wrapper = append(wrapper, uint8(len(command)-1)) // Data length
		wrapper = append(wrapper, command[1:]...)        // Data
	}
	wrapper = append(wrapper, 0x00) // Le byte

	return wrapper
}
