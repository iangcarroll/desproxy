package main

import (
	"log"
	"time"
)

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

	log.Println(initEmulation(emulator))
	for {
		command, err := receiveCommand(emulator)
		log.Println(command, err)

		time.Sleep(time.Second)
	}
}
