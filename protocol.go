// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package main

import (
	"io"
)

type Request struct {
	data string
}

func ReadRequest(r io.Reader) (req *Request, err error) {
	return nil, err
}
