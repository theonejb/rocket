/*
	Defines the constants, vars and types used by the main rocket application
*/
package main

import (
	"io"
)

const (
	uploadHandlerLocal = "local"
	uploadHandlerAws   = "aws"

	_1mb = 1 << 20
)

type serverOpts struct {
	listenAddress    string
	maxImageHandlers int
	maxImageSize     int
	uploadHandler    uploadHandler
}

type uploadHandler interface {
	name() string
	upload(string, io.Reader, int64) (string, error)
}
