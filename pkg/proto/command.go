package proto

const (
	Version = 1
)

const (
	CmdWelcome = iota
	CmdClose
	CmdSignin
	CmdSignout
	CmdConnectToMatch
	CmdMatchUpdate
)

type Command struct {
	// Header
	Version byte   // 1
	Command uint16 // 2
	ID      uint32 // 4
	// Data
	Data []byte
}
