package uploaders

import (
	"io"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

func AwsUpload(key string, data io.Reader, dataSize int64) (string, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return "", err
	}

	s3accessor := s3.New(auth, aws.APSoutheast)
	bucket := s3accessor.Bucket("jb-testbucket-1001")

	if err = bucket.PutReader(key, data, dataSize, "", "public-read"); err != nil {
		return "", err
	}

	fileUrl := bucket.URL(key)
	return fileUrl, nil
}
