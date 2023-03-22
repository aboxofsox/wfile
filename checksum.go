package wfile

import (
	"crypto/md5"
	"fmt"
	"os"
)

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
	defer f.Close()

	content := make([]byte, info.Size())
	_, err = f.Read(content)
	if err != nil {
		return [md5.Size]byte{}, err
	}

	return md5.Sum(content), nil
}
