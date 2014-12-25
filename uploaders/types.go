package uploaders

import (
	"io"
)

type UploadHandler func(key string, data io.Reader, dataSize int64) (string, error)
