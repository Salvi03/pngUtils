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
	size  uint32
	ctype string
	data  []byte
	crc   uint32
}
