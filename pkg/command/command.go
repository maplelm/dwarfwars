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
	Version CommandVersion // 8 bits
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
	bytes[0] = byte(c.Version)
	bytes[1] = byte(c.Type)
	bytes[2] = byte((c.Size & 0xFF00) >> 8)
	bytes[3] = byte(c.Size & 0xFF)
	return append(bytes, c.Data...)

}

func Unmarshal(d []byte) (*Command, error) {
	if len(d) < 4 {
		return nil, fmt.Errorf("Malformed Command (%d): %s", len(d), string(d))
	}
	return &Command{
		Version: CommandVersion(d[0]),
		Type:    CommandType(d[1]),
		Size:    ((uint16(d[2]) << 8) + uint16(d[3])),
		Data:    d[4:],
	}, nil
}

func ValidateHeader(header []byte) (msgSize uint16, cmd CommandType, err error) {
	if len(header) < 4 {
		err = fmt.Errorf("malformed header")
		return
	}
	if CommandVersion(header[0]) != CurrentVersion {
		err = fmt.Errorf("version missmatch, header version :%d, Expected Version: %d", CommandVersion(header[0]), CurrentVersion)
		return
	}

	msgSize = uint16((uint16(header[2]) << 8) + uint16(header[3]))
	cmd = CommandType(header[1])

	return
}
