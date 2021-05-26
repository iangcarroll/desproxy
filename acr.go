package main

import (
	"log"

	"github.com/ebfe/scard"
)

type EmulationRequest struct {
	Mode          byte   // 1 byte
	SensRes       []byte // 2 bytes
	Uid           []byte // 3 bytes
	SelRes        byte   // 1 byte
	Felicia       []byte // 18 bytes
	NfcId         []byte // 10 bytes
	GeneralAtr    []byte // 47 byte maximum
	HistoricalAts []byte // 48 byte maximum
}

func (e *EmulationRequest) Serialize() []byte {
	// Start off with the mode byte...
	serialized := []byte{e.Mode}

	// Add the SensRes.
	if len(e.SensRes) != 2 {
		panic("e.SensRes length invalid")
	}
	serialized = append(serialized, e.SensRes...)

	// Add the UID.
	if len(e.Uid) != 3 {
		panic("e.Uid length invalid")
	}
	serialized = append(serialized, e.Uid...)

	// Add the SelRes.
	serialized = append(serialized, e.SelRes)

	// Add the Felicia bytes.
	if len(e.Felicia) != 18 {
		panic("e.Felicia length invalid")
	}
	serialized = append(serialized, e.Felicia...)

	// Add the NfcId bytes.
	if len(e.NfcId) != 10 {
		panic("e.NfcId length invalid")
	}
	serialized = append(serialized, e.NfcId...)

	// Add the GeneralAtr length and bytes.
	serialized = append(serialized, uint8(len(e.GeneralAtr)))
	serialized = append(serialized, e.GeneralAtr...)

	// Add the HistoricalAts length and bytes.
	serialized = append(serialized, uint8(len(e.HistoricalAts)))
	serialized = append(serialized, e.HistoricalAts...)

	log.Println("Serialized EmuReq to", len(serialized), "bytes")
	return serialized
}

var (
	nxpInitiateEmulation = []byte{0xd4, 0x8c}
)

func directTransmit(card *scard.Card, command []byte) ([]byte, error) {
	directTransmit := []byte{0xFF, 0x00, 0x00, 0x00, uint8(len(command))}
	directTransmit = append(directTransmit, command...)

	return card.Control(acsControlCommand, directTransmit)
}

func initEmulation(card *scard.Card) ([]byte, error) {
	req := EmulationRequest{
		Mode:          0x00,
		SensRes:       []byte{0x44, 0x03},
		Uid:           []byte{0x98, 0x65, 0xD2},
		SelRes:        0x20,
		Felicia:       []byte{0x01, 0xfe, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xff, 0xff},
		NfcId:         []byte{0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11},
		GeneralAtr:    []byte{},
		HistoricalAts: []byte{0x80},
	}

	command := append(nxpInitiateEmulation, req.Serialize()...)
	return directTransmit(card, command)
}

func receiveCommand(card *scard.Card) ([]byte, error) {
	return directTransmit(card, []byte{0xd4, 0x86})
}
