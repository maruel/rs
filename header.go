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
	"bytes"
	"encoding/json"
	"io"
)

type rSHeader struct {
	Poly, C int
	Alpha   byte
}

func readHeader(r io.Reader) (io.Reader, *Field, int, error) {
	var hdr rSHeader
	dec := json.NewDecoder(r)
	if err := dec.Decode(&hdr); err != nil {
		return nil, nil, 0, err
	}
	b := dec.Buffered()
	p := make([]byte, 16)
	f := NewField(hdr.Poly, hdr.Alpha)
	// strip header-followed linefeeds
	for {
		n, err := b.Read(p[:cap(p)])
		if err != nil {
			return nil, nil, 0, err
		}
		p = p[:n]
		for i, c := range p {
			if c == '\n' {
				continue
			}
			p = p[i:]
			return io.MultiReader(bytes.NewReader(p), b, r), f, hdr.C, nil
		}
	}
}

func writeHeader(w io.Writer, f *Field, c int) (int, error) {
	cw := &countingWriter{Writer: w}

	hdr := rSHeader{C: c}
	hdr.Poly, hdr.Alpha = f.GetPolyAlpha()
	err := json.NewEncoder(cw).Encode(hdr)
	return int(cw.n), err
}

type countingWriter struct {
	io.Writer
	n int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.Writer.Write(p)
	cw.n += int64(n)
	return n, err
}
