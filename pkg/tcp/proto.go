package tcp

const (
	Version          byte = 1
	versionBytes     byte = 1
	commandBytes     byte = 2
	packetSizeBytes  byte = 2
	HeaderSize       byte = commandBytes + packetSizeBytes + versionBytes
	MaxDataPerPacket int  = 65535 // bytes
)
