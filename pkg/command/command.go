package command

/*
	I need to get this encoded in BigEndian and make it just make more sense but it works pretty well right now.
*/
import (
	"fmt"
	"math"
)

const CurrentVersion CommandVersion = 1

const HeaderSize int = 6 // bytes

type CommandType uint8
type CommandVersion uint8

type Command struct {
	Version CommandVersion // 8 bits
	Type    CommandType    // 8 bits
	Size    uint32         // 16 bits (max command size: 65_536 bytes)
	Data    []byte
}

func New(t CommandType, d []byte) (*Command, error) {
	if len(d) > int(math.Pow(2, 32)) {
		return nil, fmt.Errorf("command exceeds data limit")
	}
	return &Command{
		Version: CurrentVersion,
		Type:    t,
		Size:    uint32(len(d)),
		Data:    d,
	}, nil
}

func (c Command) Marshal() []byte {
	bytes := make([]byte, HeaderSize)
	bytes[0] = byte(c.Version)
	bytes[1] = byte(c.Type)
	bytes[2] = byte((c.Size & 0xFF000000) >> 24)
	bytes[3] = byte((c.Size & 0x00FF0000) >> 16)
	bytes[4] = byte((c.Size & 0x0000FF00) >> 8)
	bytes[5] = byte((c.Size & 0x000000FF))
	return append(bytes, c.Data...)

}

func Unmarshal(d []byte) (*Command, error) {
	s, t, e := ValidateHeader(d[:HeaderSize])
	if e != nil {
		return nil, e
	}
	return &Command{
		Version: CommandVersion(d[0]),
		Type:    t,
		Size:    s,
		Data:    d[HeaderSize : int(s)+HeaderSize],
	}, nil
}

func ValidateHeader(header []byte) (msgSize uint32, cmd CommandType, err error) {
	if len(header) < HeaderSize {
		err = fmt.Errorf("malformed header")
		return
	}
	if CommandVersion(header[0]) != CurrentVersion {
		err = fmt.Errorf("version missmatch, header version :%d, Expected Version: %d", CommandVersion(header[0]), CurrentVersion)
		return
	}

	msgSize = (uint32(header[2]) << 24) + (uint32(header[3]) << 16) + (uint32(header[4]) << 8) + (uint32(header[5]))
	cmd = CommandType(header[1])

	return
}
