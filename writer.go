package pngutils

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/draw"
	"image/png"
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

	err = binary.Write(file, binary.BigEndian, []byte("\x89\x50\x4e\x47\x0d\x0a\x1a\x0a"))
	if err != nil {
		return nil, err
	}

	err = writer.WriteChunk(ihdr)
	if err != nil {
		return nil, err
	}

	return writer, err
}

func (writer *ImageWriter) WriteChunk(c *Chunk) error {
	var data []byte
	var err error

	data, err = c.DataToBytes()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(writer.filename, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	defer file.Close()
	if err != nil {
		return err
	}

	err = binary.Write(file, binary.BigEndian, data)

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

func messageToLSB(message string) ([]byte, error) {
	var err error
	var result = make([]byte, len(message)*4)
	var messagebytes = []byte(message)
	var mb byte
	var i = 0
	var j = 0
	var count = 0

	for _, mbyte := range messagebytes {
		mb = mbyte
		for i < 4 {
			j = i
			for j > 0 {
				mb = mb >> 2
				j--
			}
			result[count] = mb & 0xFC

			i++
			count++
		}
	}

	return result, err
}

// WriteLSB not implemented yet
func (writer *ImageWriter) WriteLSB(message string) error {
	var err error
	var buffer *bufio.Reader
	var file *os.File
	var im image.Image
	var rect image.Rectangle
	var dst *image.NRGBA

	file, err = os.Open(writer.filename)
	if err != nil {
		return err
	}

	buffer = bufio.NewReader(file)
	im, err = png.Decode(buffer)

	rect = image.Rect(0, 0, im.Bounds().Dx(), im.Bounds().Dy())
	dst = image.NewNRGBA(rect)

	draw.Draw(dst, dst.Bounds(), im, im.Bounds().Min, draw.Src)

	return err
}
