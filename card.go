package main

import (
	"encoding/binary"
	"log"

	"github.com/ebfe/scard"
)

var (
	acsControlCommand = binary.BigEndian.Uint32([]byte{0x00, 0x31, 0x36, 0xB0})
)

func connectToCard(offset uint8, emulator bool) (*scard.Card, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}

	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, err
	}

	reader := readers[offset]
	log.Println("Connecting to reader", reader)

	protocol := scard.ProtocolUndefined
	if !emulator {
		protocol = scard.ProtocolT1
	}

	// Connect to the reader.
	card, err := ctx.Connect(reader, scard.ShareDirect, protocol)
	if err != nil {
		return nil, err
	}

	// Set the LED to ensure our control commands work.
	//_, err = card.Control(acsControlCommand, []byte{0xff, 0x00, 0x40, 0x0f, 0x04, 0x00, 0x00, 0x00, 0x00})
	//if err != nil {
	//	return nil, err
	//}

	return card, nil
}

func warmResetCard(card *scard.Card) {
	card.Reconnect(scard.ShareDirect, card.ActiveProtocol(), scard.ResetCard)
}

func coldResetCard(card *scard.Card) {
	card.Reconnect(scard.ShareDirect, card.ActiveProtocol(), scard.UnpowerCard)
}
