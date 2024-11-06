package command

/*
	I need to get this encoded in BigEndian and make it just make more sense but it works pretty well right now.


	HEADER:

	Byte 0  : Version
	byte 1-4: Client ID
	Byte 5  : Command Type
	Byte 6  : Data Format
	Byte 7-10: Size
*/
import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

const CurrentVersion CommandVersion = 1

const HeaderSize int = 11 // bytes

type CommandType uint8
type CommandVersion uint8

const (
	TypeWelcome = iota
	TypeLobbyJoinRequest
	TypeLobbyLeaveRequest
	TypeInput
	TypeStartGame
	TypeWorldData
	TypeWorldUpdate
)

const (
	FormatJSON = iota
	FormatGLOB
	FormatCSV
	FormatText
)

type Command struct {
	Version  CommandVersion `json:"Version"`   // 1 Byte
	ClientID uint32         `json:"Client ID"` // 4 Bytes
	Type     CommandType    `json:"Type"`      // 1 Byte
	Format   uint8          `json:"Format"`    // 1 Byte
	Size     uint32         `json:"Size"`      // 4 Bytes
	Data     []byte         `json:"Data"`      // Size Bytes
}

func New(id uint32, format uint8, t CommandType, d []byte) (*Command, error) {
	if len(d) > int(math.Pow(2, 32)) {
		return nil, fmt.Errorf("command exceeds data limit")
	}

	return &Command{
		Version:  CurrentVersion,
		ClientID: id,
		Type:     t,
		Format:   format,
		Size:     uint32(len(d)),
		Data:     d,
	}, nil
}

func (c Command) Send(conn net.Conn) (int64, error) {
	bytes := make([]byte, HeaderSize)
	bytes = append([]byte{}, byte(c.Version))
	bytes = binary.BigEndian.AppendUint32(bytes, c.ClientID)
	bytes = append(bytes, byte(c.Type))
	bytes = append(bytes, byte(c.Format))
	bytes = binary.BigEndian.AppendUint32(bytes, c.Size)

	var b net.Buffers = [][]byte{bytes, c.Data}

	return b.WriteTo(conn)

}

func Recieve(conn net.Conn) (*Command, error) {
	h := make([]byte, HeaderSize)
	n, err := conn.Read(h)
	if err != nil {
		return nil, err
	}
	if n != HeaderSize {
		return nil, fmt.Errorf("recieve did not get the approprate number of header bytes")
	}
	s, t, id, f, e := validateheader(h)
	if e != nil {
		return nil, e
	}
	d := make([]byte, s)
	n, err = conn.Read(d)
	if err != nil {
		return nil, err
	}
	if n != int(s) {
		return nil, fmt.Errorf("recieve did not recieve the correct number of data bytes")
	}
	return &Command{
		Version:  CommandVersion(h[0]),
		ClientID: id,
		Type:     t,
		Format:   f,
		Size:     s,
		Data:     d,
	}, nil
}

func validateheader(header []byte) (msgSize uint32, cmd CommandType, id uint32, format uint8, err error) {
	if len(header) < HeaderSize {
		err = fmt.Errorf("malformed header")
		return
	}
	if CommandVersion(header[0]) != CurrentVersion {
		err = fmt.Errorf("version missmatch, header version :%d, Expected Version: %d", CommandVersion(header[0]), CurrentVersion)
		return
	}

	id = (uint32(header[1]) << 24) + (uint32(header[2]) << 16) + (uint32(header[3]) << 8) + uint32(header[4])
	cmd = CommandType(header[5])
	format = uint8(header[6])

	msgSize = (uint32(header[7]) << 24) + (uint32(header[8]) << 16) + (uint32(header[9]) << 8) + (uint32(header[10]))

	return
}
