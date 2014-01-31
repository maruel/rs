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

// NewWriter returns a new writer with c ECC-length
func NewWriter(w io.Writer, c int) io.Writer {
	return &rSWriter{
		dataLen: maxDataLen - c,
		eccLen:  c,
		e:       NewEncoder(QR_CODE_FIELD_256, c),
		w:       w,
		block:   make([]byte, maxDataLen),
		rest:    make([]byte, 0, maxDataLen-c)}
}

func (wr *rSWriter) Write(p []byte) (int, error) {
	dataLen := wr.dataLen
	if len(p)+len(wr.rest) < dataLen {
		wr.rest = append(wr.rest, p...)
		return 0, nil
	}

	n := 0
	i := len(wr.rest)
	if i > 0 {
		copy(wr.block, wr.rest)
		wr.rest = wr.rest[:0]
	}

	k := dataLen - (len(p) + i)
	for k > 0 {
		copy(wr.block[i:], p[:k])
		p = p[k:]
		i, k = 0, dataLen-len(p)
		wr.e.Encode(wr.block[:dataLen], wr.block[dataLen:])
		if m, err := wr.w.Write(wr.block); err != nil {
			return m, err
		}
		n += dataLen
	}

	if len(p) > 0 {
		wr.rest = append(wr.rest, p...)
	}
	return n, nil
}

func (wr *rSWriter) Flush() (int, error) {
	n := len(wr.rest)
	if n == 0 {
		return 0, nil
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
	return wr.w.Write(wr.block[:n+wr.eccLen])
}

func (wr *rSWriter) Close() error {
	_, err := wr.Flush()
	return err
}
