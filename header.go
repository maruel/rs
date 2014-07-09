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
	"encoding/json"
	"io"
)

type rSHeader struct {
	Poly, C int
	Alpha   byte
}

func readHeader(r io.Reader) (*Field, int, error) {
	var hdr rSHeader
	if err := json.NewDecoder(r).Decode(&hdr); err != nil {
		return nil, 0, err
	}
	return NewField(hdr.Poly, hdr.Alpha), hdr.C, nil
}

func writeHeader(w io.Writer, f *Field, c int) (int, error) {
	cw := &countingWriter{Writer: w}

	// FIXME(tgulacsi): get poly and alpha from *Field
	err := json.NewEncoder(cw).Encode(rSHeader{Poly: 0, Alpha: 0, C: c})
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
