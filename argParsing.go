// Contains stuff related to handling uploads
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"rocket/uploaders"
)

type localUploadHandler struct {
	outDir string
}

type awsUploadHandler struct{}

func newLocalUploadHandler(outDir string) uploadHandler {
	return localUploadHandler{outDir}
}

func (uh localUploadHandler) name() string {
	return fmt.Sprintf("%s -> [%s]", uploadHandlerLocal, uh.outDir)
}

func (uh localUploadHandler) upload(key string, data io.Reader, dataSize int64) (string, error) {
	completeFileName := path.Join(uh.outDir, key)
	return uploaders.LocalUpload(completeFileName, data, dataSize)
}

func newAwsUploadHandler() uploadHandler {
	return awsUploadHandler{}
}

func (uh awsUploadHandler) name() string {
	return uploadHandlerAws
}

func (uh awsUploadHandler) upload(key string, data io.Reader, dataSize int64) (string, error) {
	return uploaders.AwsUpload(key, data, dataSize)
}

func parseCliArgs() (serverOpts, error) {
	flagBindAddress := flag.String("bind-address", "127.0.0.1:8080", "The address to listen on. Format: <IP>:<PORT>")
	flagMaxImageHandlers := flag.Int("max-image-handlers", 4, "The max. number of go-routines processing images at the same time")
	flagMaxImageSize := flag.Int("max-image-size", 10*_1mb, "Max. image size to process. Any requests with a larger image will be rejected")
	flagUploadHandlerName := flag.String("upload-handler", uploadHandlerLocal, "The upload handler to use. Valid values: local, aws")
	flagLocalOutDir := flag.String("uploader-out-dir", "/tmp", "Where to save the processed images. Only valid with the local upload handler")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of 'rocket' image processor and uploader:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var bindAddress string = *flagBindAddress
	var maxImageHandlers int = *flagMaxImageHandlers
	var maxImageSize int = *flagMaxImageSize

	var runOpts serverOpts = serverOpts{
		listenAddress:    bindAddress,
		maxImageHandlers: maxImageHandlers,
		maxImageSize:     maxImageSize,
	}

	var uploadHandlerName string = *flagUploadHandlerName
	if uploadHandlerName == uploadHandlerLocal {
		runOpts.uploadHandler = newLocalUploadHandler(*flagLocalOutDir)
	} else if uploadHandlerName == uploadHandlerAws {
		runOpts.uploadHandler = newAwsUploadHandler()
	} else {
		return runOpts, errors.New(fmt.Sprintf("Unknown uploder type '%s'", uploadHandlerName))
	}

	return runOpts, nil
}
