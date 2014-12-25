package uploaders

import (
	"errors"
	"io"
	"os"
)

func LocalUpload(key string, data io.Reader, dataSize int64) (string, error) {
	newFile, err := os.Create(key)
	if err != nil {
		return "", err
	}

	writtenBytes, err := io.Copy(newFile, data)
	if err != nil {
		return "", err
	}
	if writtenBytes != dataSize {
		return "", errors.New("Incomplete data saved. Assume corrupted. Should retry")
	}

	return key, nil
}
