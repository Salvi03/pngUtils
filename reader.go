package pngutils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
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

	c = &Chunk{}
	err = binary.Read(im.reader, binary.BigEndian, &c.Size)

	var t = make([]byte, 4)
	err = binary.Read(im.reader, binary.BigEndian, &t)

	c.Data = make([]byte, c.Size)
	err = binary.Read(im.reader, binary.BigEndian, &c.Data)

	c.Ctype = string(t)
	err = binary.Read(im.reader, binary.BigEndian, &c.Crc)

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

	c = &Chunk{}
	c.Ctype = ""
	for c.Ctype != "IEND" {
		c, err = im.ReadChunk()
		if err != nil {
			break
		}

		cs = append(cs, c)
	}

	return cs, err
}

type colors struct {
	red   bool
	green bool
	blue  bool
	x     int
	y     int
}

func getLSBContent(size uint32, img *image.NRGBA, col *colors) []byte {
	var pixel color.NRGBA
	var resultLSB = make([]byte, 4*size)
	var I = uint32(0)
	var i int
	var j int
	var k int
	var result = make([]byte, size)

	index := 0

	for I < size*4 {
		pixel = img.NRGBAAt(col.x, col.y)

		if !col.red {
			resultLSB[I] = pixel.R & 0x03
			I++
			col.red = true
		}

		if I >= size*4 {
			break
		}

		if !col.green {
			resultLSB[I] = pixel.G & 0x03
			I++
			col.green = true
		}

		if I >= size*4 {
			break
		}

		if !col.blue {
			resultLSB[I] = pixel.B & 0x03
			I++
			col.blue = true
		}

		if I >= size*4 {
			break
		}

		col.red = false
		col.green = false
		col.blue = false

		if col.x < img.Bounds().Dx() {
			col.x++
		} else {
			col.x = 0
			col.y++
		}
	}

	var offset = uint32(0)

	var char byte
	var charLSB = make([]byte, 4)
	I = 0
	index = 0

	for offset < size*4 {
		I = 0
		for I < 4 {
			charLSB[I] = resultLSB[offset+I]
			I++
		}

		i = 3
		char = 0x00
		j = 0
		k = 0

		for i > 0 {
			k = 0
			for k < i {
				charLSB[i] = charLSB[i] << 2
				k++
			}

			i--
			j++
		}

		for _, b := range charLSB {
			char += b
		}

		result[index] = char
		offset += 4
		index++
	}

	return result
}

func getLSBMessage(img *image.NRGBA) ([]byte, error) {
	var result []byte
	var err error

	// var index = 0
	col := &colors{
		red:   false,
		blue:  false,
		green: false,
		x:     0,
		y:     0,
	}

	bsize := getLSBContent(4, img, col)

	size := binary.BigEndian.Uint32(bsize)
	result = getLSBContent(size, img, col)

	return result, err
}

func ReadLSBMessage(filename string) ([]byte, error) {
	var err error
	var file *os.File
	var buffer *bufio.Reader
	var img image.Image
	var rgba *image.NRGBA
	var rect image.Rectangle

	var result []byte

	file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	buffer = bufio.NewReader(file)
	img, err = png.Decode(buffer)
	if err != nil {
		return nil, err
	}
	rect = image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())
	rgba = image.NewNRGBA(rect)

	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	result, err = getLSBMessage(rgba)

	return result, err
}
