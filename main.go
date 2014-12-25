package main

import (
	"fmt"
	"os"
)

func main() {
	runOpts, err := parseCliArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing command line args\n\t%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Running with Options:\n\tListening Addr: %s\n\tMax. Image Handlers: %d\n\tMax. Image Size: %0.3fMb\n\tUpload Handler: %s\n",
		runOpts.listenAddress, runOpts.maxImageHandlers, float32(runOpts.maxImageSize)/_1mb, runOpts.uploadHandler.name())

	startServer(runOpts)
}
