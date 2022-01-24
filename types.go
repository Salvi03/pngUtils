package pngUtils

import (
	"bytes"
	"os"
)

type imageReader struct {
	filename string
	reader   *bytes.Reader
}

type imageWriter struct {
	filename string
	file     *os.File
}

type chunk struct {
	length uint32
	ctype  string
	data   []byte
	crc    uint32
}
