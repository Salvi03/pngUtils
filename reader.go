package pngutils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

func (im *imageReader) initialize() (*chunk, error) {
	var file *os.File
	var err error
	var buffer *bufio.Reader
	var stats os.FileInfo
	var ihdr *chunk

	file, err = os.Open(im.filename)
	buffer = bufio.NewReader(file)

	stats, err = file.Stat()
	content := make([]byte, stats.Size())

	_, err = buffer.Read(content)
	im.reader = bytes.NewReader(content)
	if err != nil {
		return nil, err
	}

	ihdr, err = im.validate()
	return ihdr, err
}

func (im *imageReader) validate() (*chunk, error) {
	var header []byte
	var h = make([]byte, 8)

	var buf uint64
	var err error

	header = make([]byte, 8)
	copy(header, "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a")

	err = binary.Read(im.reader, binary.BigEndian, &buf)
	binary.BigEndian.PutUint64(h, buf)

	result := bytes.Compare(header, h)
	if result != 0 {
		err = errors.New("this is not a valid PNG file")
	}

	c, err1 := im.readChunk()
	if err1 != nil {
		return nil, err1
	}

	return c, err
}

func (im *imageReader) readChunk() (*chunk, error) {
	var c *chunk
	var err error

	err = binary.Read(im.reader, binary.BigEndian, &c.length)

	var t = make([]byte, 4)
	err = binary.Read(im.reader, binary.BigEndian, &t)

	c.ctype = string(t)
	err = binary.Read(im.reader, binary.BigEndian, &c.crc)

	return c, err
}

func (im *imageReader) resetReader() (*chunk, error) {
	var err error
	var ihdr *chunk

	_, err = im.reader.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	ihdr, err = im.validate()
	return ihdr, err
}

func (im *imageReader) readChunkPosition(n int) (*chunk, error) {
	var c *chunk
	var err error
	var i = 0

	_, err = im.reader.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = im.validate()
	if err != nil {
		return nil, err
	}

	for i < n {
		c, err = im.readChunk()
		if err != nil {
			break
		}

		i++
	}

	_, err = im.resetReader()
	if err != nil {
		return nil, err
	}

	return c, err
}

func (im *imageReader) readNChunks(n int) ([]*chunk, error) {
	var err error
	var chunks []*chunk
	var c *chunk
	var i = 0

	for i < n {
		c, err = im.readChunk()
		chunks = append(chunks, c)
		if err != nil {
			break
		}

		i++
	}

	return chunks, err
}

func (im *imageReader) readChunksTillTheEnd() ([]*chunk, error) {
	var c *chunk
	var cs []*chunk

	var err error

	c.ctype = ""
	for c.ctype != "IEND" {
		c, err = im.readChunk()
		if err != nil {
			break
		}

		cs = append(cs, c)
	}

	return cs, err
}
