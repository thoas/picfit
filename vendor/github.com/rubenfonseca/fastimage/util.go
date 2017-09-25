package fastimage

import (
	"bytes"
	"encoding/binary"
)

func readUint16(buffer []byte) (result uint16) {
	reader := bytes.NewReader(buffer)
	binary.Read(reader, binary.BigEndian, &result)
	return
}

func readULint16(buffer []byte) (result uint16) {
	reader := bytes.NewReader(buffer)
	binary.Read(reader, binary.LittleEndian, &result)
	return
}
