package pngutils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func InitializeImageReader(filename string) (*ImageReader, *Chunk, error) {
	var file *os.File
	var err error
	var buffer *bufio.Reader
	var stats os.FileInfo
	var ihdr *Chunk
	var im *ImageReader

	im = &ImageReader{}
	im.filename = filename

	file, err = os.Open(im.filename)
	buffer = bufio.NewReader(file)

	stats, err = file.Stat()
	content := make([]byte, stats.Size())

	_, err = buffer.Read(content)
	im.reader = bytes.NewReader(content)
	if err != nil {
		return nil, nil, err
	}

	ihdr, err = im.validate()
	return im, ihdr, err
}

func (im *ImageReader) validate() (*Chunk, error) {
	var header []byte
	var h = make([]byte, 8)

	var buf uint64
	var err error

	header = make([]byte, 8)
	copy(header, "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a")
	fmt.Println("diao")

	err = binary.Read(im.reader, binary.BigEndian, &buf)
	binary.BigEndian.PutUint64(h, buf)

	result := bytes.Compare(header, h)
	if result != 0 {
		err = errors.New("this is not a valid PNG file")
	}

	c, err1 := im.ReadChunk()
	if err1 != nil {
		return nil, err1
	}

	return c, err
}

func (im *ImageReader) ReadChunk() (*Chunk, error) {
	var c *Chunk
	var err error

	err = binary.Read(im.reader, binary.BigEndian, &c.Length)

	var t = make([]byte, 4)
	err = binary.Read(im.reader, binary.BigEndian, &t)

	c.ctype = string(t)
	err = binary.Read(im.reader, binary.BigEndian, &c.crc)

	return c, err
}

func (im *ImageReader) ResetReader() (*Chunk, error) {
	var err error
	var ihdr *Chunk

	_, err = im.reader.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	ihdr, err = im.validate()
	return ihdr, err
}

func (im *ImageReader) ReadChunkPosition(n int) (*Chunk, error) {
	var c *Chunk
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
		c, err = im.ReadChunk()
		if err != nil {
			break
		}

		i++
	}

	_, err = im.ResetReader()
	if err != nil {
		return nil, err
	}

	return c, err
}

func (im *ImageReader) ReadNChunks(n int) ([]*Chunk, error) {
	var err error
	var chunks []*Chunk
	var c *Chunk
	var i = 0

	for i < n {
		c, err = im.ReadChunk()
		chunks = append(chunks, c)
		if err != nil {
			break
		}

		i++
	}

	return chunks, err
}

func (im *ImageReader) ReadChunksTillTheEnd() ([]*Chunk, error) {
	var c *Chunk
	var cs []*Chunk

	var err error

	c.ctype = ""
	for c.ctype != "IEND" {
		c, err = im.ReadChunk()
		if err != nil {
			break
		}

		cs = append(cs, c)
	}

	return cs, err
}
