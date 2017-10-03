package tabconv

import (
	"encoding/binary"
	"os"

	"github.com/pkg/errors"
)

const (
	fileHeaderSignature = 0x04034b50
	fileHeaderLen       = 30
)

var ErrFormat = errors.New("invalid xlsx file format")

// checkXlsxFile checks if file is ZIP archive
func checkXlsxFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, fileHeaderLen)
	if _, err := f.Read(buf); err != nil {
		return err
	}

	v := binary.LittleEndian.Uint32(buf)
	buf = (buf)[4:]

	if v != fileHeaderSignature {
		return ErrFormat
	}

	return nil
}
