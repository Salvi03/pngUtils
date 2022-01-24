package pngUtils

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

func (c *chunk) dataToBytes() ([]byte, error) {
	var data []byte
	var err error

	data = make([]byte, 4)
	binary.BigEndian.PutUint32(data, c.length)

	data = append(data, []byte(c.ctype)...)
	data = append(data, c.data...)

	crc := make([]byte, 4)
	binary.BigEndian.PutUint32(crc, c.crc)
	data = append(data, crc...)

	return data, err
}

func newChunk(ctype string, data string) (*chunk, error) {
	var c *chunk
	c = &chunk{ctype: ctype, data: []byte(data)}

	c.length = uint32(len(data))
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, c.ctype)
	err = binary.Write(buffer, binary.BigEndian, c.data)

	c.crc = crc32.ChecksumIEEE(buffer.Bytes())

	return c, err
}
