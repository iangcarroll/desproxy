package main

import "github.com/ebfe/scard"

var (
	nxpInitiateEmulation = []byte{0xd4, 0x8c}
)

func directTransmit(card *scard.Card, command []byte) ([]byte, error) {
	directTransmit := []byte{0xFF, 0x00, 0x00, 0x00, uint8(len(command))}
	directTransmit = append(directTransmit, command...)

	return card.Transmit(directTransmit)
}

func initEmulation(card *scard.Card) {
	nxpInitiateEmulation := append(nxpInitiateEmulation, 0x00)
	directTransmit(card, nxpInitiateEmulation)
}
