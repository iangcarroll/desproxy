package main

import "time"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Establish the emulator.
	emulator, err := connectToCard(0)
	check(err)

	// Establish the reader.
	_, err = connectToCard(1)
	check(err)

	for {
		directTransmit(emulator, []byte{})
		time.Sleep(time.Millisecond * 500)
	}
}
