package tcp

import (
	"encoding/binary"
	"fmt"
	"log"
)

/*
Byte 0  : Version
Byte 1-2: Length data - header
Byte 3-4: Command
Byte 5+ : Data

General Data Structure (json for visual but not sure if it will be actual json)
{
	"Session ID": "123abc",
	"UID": "Luke123vddk",
	"Session Secret: "asdfwefcb32435fasdaf243",
	"Command Relevent Data": [...]
}

*/

type Command struct {
	Command uint16
	Length  int
	Header  []byte
	Data    []byte
}

type WelcomeCmd func() *Command

func (t *Command) UnmarshalBinary(b []byte) error {
	// Checking if the incoming command version matches the expected version
	if b[1] != Version {
		return fmt.Errorf("Version mismatch, expected %d got %d", Version, b[0])
	}
	b = b[1:]
	// Getting the header
	t.Header = b[:HeaderSize]
	b = b[HeaderSize:]
	t.Command = binary.BigEndian.Uint16(t.Header[:commandBytes])
	t.Length = len(t.Header) + 2 + len(b)
	t.Data = b

	return nil
}

func (t *Command) MarshalBinary() (b []byte, err error) {

	// Length of data to be sent
	length := len(t.Data)
	if length > MaxDataPerPacket {
		log.Printf("Warning: TCP Command data exceeds Max (%d): %d", MaxDataPerPacket, length)
		t.Data = t.Data[:MaxDataPerPacket]
		length = MaxDataPerPacket
	}
	lengthData := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthData, uint16(length))

	// Command for the header
	commandData := make([]byte, commandBytes)
	binary.BigEndian.PutUint16(commandData, t.Command)

	// Crafting the packet
	b = make([]byte, 0, int(versionBytes)+int(commandBytes)+int(packetSizeBytes)+int(length))
	b = append(b, Version)
	b = append(b, commandData...)
	b = append(b, lengthData...)
	for i := range t.Data {
		b[int(HeaderSize)+i] = t.Data[i]
	}
	return
}
