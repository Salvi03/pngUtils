package pngutils

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/color"
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

func messageToLSB(message []byte) ([]byte, error) {
	var err error
	var result = make([]byte, len(message)*4)
	var mb byte
	var i = 0
	var j = 0
	var count = 0

	for _, mbyte := range message {
		i = 0
		for i < 4 {
			mb = mbyte
			j = i
			for j > 0 {
				mb = mb >> 2
				j--
			}
			result[count] = mb & 0x03

			i++
			count++
		}
	}

	return result, err
}

func writeMessage(rgba *image.NRGBA, LSBMessage []byte) {
	var pixel color.NRGBA
	var x = 0
	var y = 0
	var index = 0

	var size = uint32(len(LSBMessage))
	var bsize = make([]byte, 4)
	binary.BigEndian.PutUint32(bsize, size)

	LSBSize, _ := messageToLSB(bsize)
	LSBMessage = append(LSBSize, LSBMessage...)

	var red *byte
	var green *byte
	var blue *byte

	for index < len(LSBMessage) {
		pixel = rgba.NRGBAAt(x, y)

		red = &pixel.R
		green = &pixel.G
		blue = &pixel.B

		*red = *red & 0xFC
		*red += LSBMessage[index]
		index++
		if index >= len(LSBMessage) {
			rgba.SetNRGBA(x, y, pixel)
			continue
		}

		*green = *green & 0xFC
		*green += LSBMessage[index]
		index++
		if index >= len(LSBMessage) {
			rgba.SetNRGBA(x, y, pixel)
			continue
		}

		*blue = *blue & 0xFC
		*blue += LSBMessage[index]
		index++
		if index >= len(LSBMessage) {
			rgba.SetNRGBA(x, y, pixel)
			continue
		}

		rgba.SetNRGBA(x, y, pixel)
		if x < rgba.Bounds().Dx() {
			x++
		} else {
			x = 0
			y++
		}
	}
}

// WriteLSB writes your message in the two least significant bits of every pixel
func WriteLSB(infile string, outfile string, message string) error {
	var err error
	var buffer *bufio.Reader
	var file *os.File
	var im image.Image
	var rect image.Rectangle
	var dst *image.NRGBA
	var LSB []byte
	var out *os.File

	file, err = os.Open(infile)
	if err != nil {
		return err
	}

	buffer = bufio.NewReader(file)
	im, err = png.Decode(buffer)
	if err != nil {
		return err
	}

	rect = image.Rect(0, 0, im.Bounds().Dx(), im.Bounds().Dy())
	dst = image.NewNRGBA(rect)

	draw.Draw(dst, dst.Bounds(), im, im.Bounds().Min, draw.Src)
	LSB, err = messageToLSB([]byte(message))
	if err != nil {
		return err
	}

	writeMessage(dst, LSB)
	out, err = os.Create(outfile)
	if err != nil {
		return err
	}

	err = png.Encode(out, dst)
	out.Close()

	return err
}
