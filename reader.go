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

type rSReader struct {
	dataLen, eccLen int
	d               Decoder
	r               io.Reader
	block, rest     []byte
	corr            int64 //corrected byte count
}

// NewInterleavedReader returns a reader correcting on-the-fly.
//
// This reads the field and the c parameter from the data header.
func NewInterleavedReader(r io.Reader) (io.Reader, error) {
	r, f, c, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	return NewInterleavedReaderField(r, f, c), nil
}

// NewInterleavedReaderField returns a reader correcting on-the-fly
func NewInterleavedReaderField(r io.Reader, f *Field, c int) io.Reader {
	return &rSReader{
		dataLen: maxDataLen - c,
		eccLen:  c,
		d:       NewDecoder(f),
		r:       r,
		block:   make([]byte, maxDataLen),
	}
}

func (rdr *rSReader) Read(p []byte) (int, error) {
	if len(rdr.rest) == 0 {
		// decode next block
		dataLen := rdr.dataLen
		rdr.block = rdr.block[:cap(rdr.block)]
		n, err := io.ReadFull(rdr.r, rdr.block)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				return n, err
			}
			if n == 0 {
				return 0, io.EOF
			}
			// padding: move ecc code to the and and fill with zeroes inbetween
			copy(rdr.block[dataLen:], rdr.block[n-rdr.eccLen:])
			for i := n; i < dataLen; i++ {
				rdr.block[i] = 0
			}
		}
		corr, err := rdr.d.Decode(rdr.block[:dataLen], rdr.block[dataLen:])
		if corr > 0 {
			rdr.corr += int64(corr)
		}
		rdr.rest = rdr.block[:n-rdr.eccLen]
	}
	if len(rdr.rest) > 0 {
		n := len(p)
		if n > len(rdr.rest) {
			n = len(rdr.rest)
		}
		copy(p, rdr.rest[:n])
		rdr.rest = rdr.rest[n:]
		return n, nil
	}
	return 0, io.EOF
}
