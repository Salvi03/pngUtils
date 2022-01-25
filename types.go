package pngutils

import (
	"bytes"
	"os"
)

type ImageReader struct {
	filename string
	reader   *bytes.Reader
}

type ImageWriter struct {
	filename string
	file     *os.File
}

type Chunk struct {
	Size  uint32
	Ctype string
	Data  []byte
	Crc   uint32
}
