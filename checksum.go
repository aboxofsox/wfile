package wfile

import (
	"crypto/md5"
	"fmt"
	"os"
)

const (
	MaxFileSize = 1 << 12
)

// Checksum calculates the MD5 checksum of the file at the specified path. If the file does not exist or is too
// large to process, an error is returned. Otherwise, the function returns the checksum as a byte array and a nil error.
func Checksum(path string) ([md5.Size]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return [md5.Size]byte{}, err
	}

	if info.Size() > MaxFileSize {
		return [md5.Size]byte{}, fmt.Errorf("max file size exceeded: %s", path)
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return [md5.Size]byte{}, nil
	}
	if err != nil {
		return [md5.Size]byte{}, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	content := make([]byte, info.Size())
	_, err = f.Read(content)
	if err != nil {
		return [md5.Size]byte{}, err
	}

	return md5.Sum(content), nil
}
