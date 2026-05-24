package relay

import(
	"fmt"
	"encoding/binary"
)

type InputPacket struct {
	Tick uint64
	PlayerID uint8
	Keys uint8
	Checksum uint32
}
const PacketSize = 14
//serialize the packet to bytes
func Serialize(p InputPacket) []byte {
	buf := make([]byte, PacketSize)
	binary.BigEndian.PutUint64(buf[0:8], p.Tick)
	buf[8] = p.PlayerID
	buf[9] = p.Keys
	binary.BigEndian.PutUint32(buf[10:14], p.Checksum)
	return buf
}
//deserialize bytes to packet
func Deserialize(data []byte) (InputPacket, error) {
	if len(data) < PacketSize {
		return InputPacket{}, fmt.Errorf("data too short to deserialize")
	}
	p := InputPacket{
		Tick: binary.BigEndian.Uint64(data[0:8]),
		PlayerID: data[8],
		Keys: data[9],
		Checksum: binary.BigEndian.Uint32(data[10:14]),
	}
	return p, nil
}

