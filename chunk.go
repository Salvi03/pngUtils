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
	binary.BigEndian.PutUint32(data, c.Size)

	data = append(data, []byte(c.Ctype)...)
	data = append(data, c.Data...)

	crc := make([]byte, 4)
	binary.BigEndian.PutUint32(crc, c.Crc)
	data = append(data, crc...)

	return data, err
}

func NewChunk(ctype string, data string) (*Chunk, error) {
	var c *Chunk
	var err error

	c = &Chunk{Ctype: ctype, Data: []byte(data)}

	c.Size = uint32(len(data))
	buffer := new(bytes.Buffer)

	c.Ctype = ctype
	c.Data = []byte(data)

	c.Crc = crc32.ChecksumIEEE(buffer.Bytes())

	return c, err
}
