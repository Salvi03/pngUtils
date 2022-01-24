package pngUtils

import (
	"encoding/binary"
	"os"
)

func initializeWriter(filename string, ihdr *chunk) (*imageWriter, error) {
	var file *os.File
	var err error

	file, err = os.Create(filename)
	if err != nil {
		return nil, err
	}

	var writer *imageWriter
	writer = &imageWriter{
		filename: filename,
		file:     file,
	}

	err = binary.Write(file, binary.BigEndian, "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a")
	if err != nil {
		return nil, err
	}

	err = writer.writeChunk(ihdr)
	return writer, err
}

func (writer *imageWriter) writeChunk(c *chunk) error {
	var data []byte
	var err error

	data, err = c.dataToBytes()
	if err != nil {
		return err
	}

	err = binary.Write(writer.file, binary.BigEndian, data)
	return err
}

func (writer *imageWriter) writeChunks(cs []*chunk) error {
	var err error

	for _, c := range cs {
		err = writer.writeChunk(c)
		if err != nil {
			break
		}
	}

	return err
}
