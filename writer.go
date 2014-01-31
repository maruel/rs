/*
Copyright 2014 Tamás Gulácsi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law
or agreed to in writing, software distributed under the License is
distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing
permissions and limitations under the License.
*/

package rs

import (
	"io"
)

const maxDataLen = 255

type rSWriter struct {
	dataLen, eccLen int
	e               Encoder
	w               io.Writer
	block           []byte
	rest            []byte
}

// NewWriter returns a new writer with f field and c ECC-length
func NewWriter(w io.Writer, f *Field, c int) io.WriteCloser {
	return &rSWriter{
		dataLen: maxDataLen - c,
		eccLen:  c,
		e:       NewEncoder(f, c),
		w:       w,
		block:   make([]byte, maxDataLen),
		rest:    make([]byte, 0, maxDataLen-c)}
}

func (wr *rSWriter) Write(p []byte) (int, error) {
	lenp := len(p)
	dataLen := wr.dataLen
	if len(p)+len(wr.rest) < dataLen {
		wr.rest = append(wr.rest, p...)
		return lenp, nil
	}

	// send rest
	ecc := make([]byte, wr.eccLen)
	for i := len(wr.rest); i >= dataLen; i = len(wr.rest) {
		if _, err := wr.w.Write(wr.rest); err != nil {
			return lenp, err
		}
		wr.e.Encode(wr.rest[:dataLen], ecc)
		if _, err := wr.w.Write(ecc); err != nil {
			return lenp, err
		}
		wr.rest = wr.rest[dataLen:]
	}

	// send rest of rest plus begin of p
	i := len(wr.rest)
	if i > 0 {
		copy(wr.block, wr.rest)
		wr.rest = wr.rest[:0]
		k := dataLen - (len(p) + i)
		if k > 0 {
			wr.rest = append(wr.rest, p...)
			return lenp, nil
		}
		copy(wr.block[i:], p[:dataLen-i])
		wr.e.Encode(wr.block[:dataLen], wr.block[dataLen:])
		if _, err := wr.w.Write(wr.block); err != nil {
			return lenp, err
		}
		p = p[dataLen-i:]
	}

	// walk on p
	for i := len(p); i >= dataLen; i = len(p) {
		if _, err := wr.w.Write(p[:dataLen]); err != nil {
			return lenp, err
		}
		wr.e.Encode(p[:dataLen], ecc)
		if _, err := wr.w.Write(ecc); err != nil {
			return lenp, err
		}
		p = p[dataLen:]
	}

	if len(p) > 0 {
		wr.rest = append(wr.rest, p...)
	}
	return lenp, nil
}

func (wr *rSWriter) Flush() error {
	n := len(wr.rest)
	if n == 0 {
		return nil
	}
	copy(wr.block, wr.rest)
	dataLen := wr.dataLen
	// shortening means fillup with zeroes, but don't transmit them
	for i := n; i < dataLen; i++ {
		wr.block[i] = 0
	}
	wr.e.Encode(wr.block[:dataLen], wr.block[dataLen:])
	wr.rest = wr.rest[:0]
	// don't transmit means cut
	copy(wr.block[n:], wr.block[dataLen:])
	_, err := wr.w.Write(wr.block[:n+wr.eccLen])
	return err
}

func (wr *rSWriter) Close() error {
	return wr.Flush()
}
