// vim:set ts=2 sw=2 et ai ft=go:
// Copyright (c) 2012 Toby DiPasquale. See accompanying LICENSE file for
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

