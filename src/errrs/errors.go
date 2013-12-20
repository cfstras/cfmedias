// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modified by cfstras to implement json.Marshaler

// Package errrs implements functions to manipulate errors.
package errrs

import (
	"encoding/json"
)

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func (e *errorString) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.s)
	//return []byte(e.s), nil
}
