package main

import (
	"net/http"
	"time"
)

const (
	serverReadTimeout             = 10
	serverWriteTimeout            = 10
	maxTimeToWaitForFreeProcessor = 5

	incomingDataReadBufferSize = _1mb
)

type imageProcessingHandler struct {
	runOpts               serverOpts
	freeFileReaderBuffers chan []byte // A free list of buffers that are used to read incomming file data
	freeFileDataBuffers   chan []byte
	rateLimiter           chan bool // We read from this when we start processing. Write back when done. Thus limiting the max number of req processed
}

func startServer(runOpts serverOpts) {
	h := &imageProcessingHandler{
		runOpts:               runOpts,
		freeFileReaderBuffers: make(chan []byte, runOpts.maxImageHandlers),
		freeFileDataBuffers:   make(chan []byte, runOpts.maxImageHandlers),
		rateLimiter:           make(chan bool, runOpts.maxImageHandlers),
	}

	// Fill free lists
	for i := 0; i < runOpts.maxImageHandlers; i++ {
		buffer := make([]byte, incomingDataReadBufferSize)
		select {
		case h.freeFileReaderBuffers <- buffer:
			continue
		default:
			break
		}
	}
	for i := 0; i < runOpts.maxImageHandlers; i++ {
		buffer := make([]byte, runOpts.maxImageSize)
		select {
		case h.freeFileDataBuffers <- buffer:
			continue
		default:
			break
		}
	}
	for i := 0; i < runOpts.maxImageHandlers; i++ {
		select {
		case h.rateLimiter <- true:
			continue
		default:
			break
		}
	}

	server := http.Server{
		Addr:         runOpts.listenAddress,
		Handler:      h,
		ReadTimeout:  time.Duration(serverReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverWriteTimeout) * time.Second,
	}
	server.ListenAndServe()
}

func (h *imageProcessingHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Make sure we don't exceed our maxImageHandlers limit
	select {
	case <-h.rateLimiter:
	case <-time.After(time.Duration(maxTimeToWaitForFreeProcessor) * time.Second):
		resp.WriteHeader(http.StatusServiceUnavailable)
		resp.Write([]byte("No image processors available right now. Please try later"))
		return
	}

	err := h.parseRequestParams(req)
	if err != nil {
		writeErrorToResponse(resp, err)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Test"))

	h.rateLimiter <- true // Signal that we've complete one operation
}

func (h *imageProcessingHandler) parseRequestParams(req *http.Request) error {
	_, err := req.MultipartReader()
	if err != nil {
		return err
	}

	return nil
}

func writeErrorToResponse(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(err.Error()))
}
