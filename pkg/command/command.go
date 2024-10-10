package command

/*
	I need to get this encoded in BigEndian and make it just make more sense but it works pretty well right now.
*/
import (
	"fmt"
)

const CurrentVersion CommandVersion = 1
const HeaderSize int = 32 // bits

type CommandType uint8
type CommandVersion uint8

type Command struct {
	Version CommandVersion // 4 bits
	Type    CommandType    // 8 bits
	Size    uint16         // 16 bits (max command size: 65_536 bytes)
	Data    []byte
}

func New(t CommandType, d []byte) (*Command, error) {
	if len(d) > 65_535 {
		return nil, fmt.Errorf("command exceeds data limit")
	}
	return &Command{
		Version: CurrentVersion,
		Type:    t,
		Size:    uint16(len(d)),
		Data:    d,
	}, nil
}

func (c Command) Marshal() []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(c.Version<<4) + byte(c.Type&0x0F)
	bytes[1] = byte(c.Type&0xF0) + byte(c.Size>>12)
	bytes[2] = byte((c.Size & 0xFF0) >> 4)
	bytes[3] = byte((c.Size & 0x0F) << 4)
	return append(bytes, c.Data...)

}

func Unmarshal(d []byte) (*Command, error) {
	if len(d) < 4 {
		return nil, fmt.Errorf("Malformed Command (%d): %s", len(d), string(d))
	}
	return &Command{
		Version: CommandVersion((d[0] & 0xF0) >> 4),
		Type:    CommandType(d[0]&0x0F) + CommandType((d[1] & 0xF0)),
		Size:    (uint16((d[1] & 0x0F)) << 12) + (uint16(d[2]) << 4) + uint16((d[3]&0xF0)>>4),
		Data:    d[4:],
	}, nil
}
