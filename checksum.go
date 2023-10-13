package wfile

import (
	"crypto/md5"
	"os"
)

// checksum calculates the MD5 checksum of the file at the specified Path. If the file does not exist or is too
// large to process, an Error is returned. Otherwise, the function returns the checksum as a byte array and a nil Error.
func checksum(path string) ([md5.Size]byte, error) {
	var sum [md5.Size]byte

	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return sum, nil
		}
		return sum, err
	}

	sum = md5.Sum(fileData)
	return sum, nil
}
