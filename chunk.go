package pngutils

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

func (c *Chunk) DataToBytes() ([]byte, error) {
	var data []byte
	var err error

	data = make([]byte, 4)
	binary.BigEndian.PutUint32(data, c.size)

	data = append(data, []byte(c.ctype)...)
	data = append(data, c.data...)

	crc := make([]byte, 4)
	binary.BigEndian.PutUint32(crc, c.crc)
	data = append(data, crc...)

	return data, err
}

func NewChunk(ctype string, data string) (*Chunk, error) {
	var c *Chunk
	c = &Chunk{ctype: ctype, data: []byte(data)}

	c.size = uint32(len(data))
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, c.ctype)
	err = binary.Write(buffer, binary.BigEndian, c.data)

	c.crc = crc32.ChecksumIEEE(buffer.Bytes())

	return c, err
}
