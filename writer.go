package pngutils

import (
	"encoding/binary"
	"os"
)

func InitializeWriter(filename string, ihdr *Chunk) (*ImageWriter, error) {
	var file *os.File
	var err error

	file, err = os.Create(filename)
	if err != nil {
		return nil, err
	}

	var writer *ImageWriter
	writer = &ImageWriter{
		filename: filename,
		file:     file,
	}

	err = binary.Write(file, binary.BigEndian, "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a")
	if err != nil {
		return nil, err
	}

	err = writer.WriteChunk(ihdr)
	return writer, err
}

func (writer *ImageWriter) WriteChunk(c *Chunk) error {
	var data []byte
	var err error

	data, err = c.DataToBytes()
	if err != nil {
		return err
	}

	err = binary.Write(writer.file, binary.BigEndian, data)
	return err
}

func (writer *ImageWriter) WriteChunks(cs []*Chunk) error {
	var err error

	for _, c := range cs {
		err = writer.WriteChunk(c)
		if err != nil {
			break
		}
	}

	return err
}
