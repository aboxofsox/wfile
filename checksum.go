package wfile

import (
	"crypto/md5"
	"os"
)

func Checksum(path string) ([md5.Size]byte, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return [md5.Size]byte{}, nil
	}
	if err != nil {
		return [md5.Size]byte{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return [md5.Size]byte{}, err
	}

	content := make([]byte, info.Size())
	_, err = f.Read(content)
	if err != nil {
		return [md5.Size]byte{}, err
	}

	return md5.Sum(content), nil
}
