package gologix

import (
	"bytes"
	"encoding/binary"
)

// multiReadReplyBody is one entry of a MultipleService reply carrying a Read
// response of the given CIP type. typeInfo is the byte after the type code
// (0x02 for a struct, 0x00 for atomic types).
func multiReadReplyBody(typ CIPType, typeInfo byte, data []byte) []byte {
	var b bytes.Buffer
	b.WriteByte(byte(CIPService_Read.AsResponse()))
	b.WriteByte(0)                                       // reserved
	_ = binary.Write(&b, binary.LittleEndian, uint16(0)) // general status + extended size
	b.WriteByte(byte(typ))
	b.WriteByte(typeInfo)
	b.Write(data)
	return b.Bytes()
}

// replyMultiRead sends a MultipleService reply whose per-service bodies are
// given. Offsets are relative to the reply-count field, as readList expects.
func (fs *fakeCIPServer) replyMultiRead(seq uint16, bodies [][]byte) {
	fs.t.Helper()
	var inner bytes.Buffer
	_ = binary.Write(&inner, binary.LittleEndian, uint16(len(bodies)))
	offset := 2 + 2*len(bodies)
	for _, body := range bodies {
		_ = binary.Write(&inner, binary.LittleEndian, uint16(offset))
		offset += len(body)
	}
	for _, body := range bodies {
		inner.Write(body)
	}
	fs.sendFrame(cipCommandSendUnitData, buildConnectedReply(seq, 0x00112233, CIPService_MultipleService, 0x00, inner.Bytes()))
}
