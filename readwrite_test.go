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
	"crypto/rand"
	"crypto/sha1"
	"io"
	"testing"
)

const eccLen = 32

var rounds = 32

func TestReadWrite(t *testing.T) {
	randBytes := make([]byte, maxDataLen-eccLen-1)
	n, err := rand.Read(randBytes)
	if err != nil {
		t.Fatalf("error reading rand bytes: %v", err)
	}
	randBytes = randBytes[:n]

	for i := 0; i < rounds && i < len(randBytes); i++ {
		encrypted := bytes.NewBuffer(nil)
		wr, err := NewInterleavedWriter(encrypted, QR_CODE_FIELD_256, eccLen)
		if err != nil {
			t.Errorf("NewInterleavedWriter: %v", err)
		}
		t.Logf("header: %s", encrypted.String())
		headerLen := encrypted.Len()
		hsh := sha1.New()
		var written int64

		w := io.MultiWriter(wr, hsh)
		for j := 0; j <= i; j++ {
			if n, err = w.Write(randBytes); err != nil {
				t.Fatalf("error writing: %v", err)
			}
			written += int64(n)
			if i > 0 {
				if n, err = w.Write(randBytes[:i]); err != nil {
					t.Fatalf("error writing: %v", err)
				}
				written += int64(n)
			}
		}
		if err = wr.Close(); err != nil {
			t.Fatalf("error closing: %v", err)
		}

		// check
		origHash := hsh.Sum(nil)
		t.Logf("written %d raw bytes, %d header, %d encrypted, hash is %x",
			written, headerLen, encrypted.Len()-headerLen, origHash)

		hsh = sha1.New()
		rdr, err := NewInterleavedReader(bytes.NewReader(encrypted.Bytes()))
		if err != nil {
			t.Errorf("NewInterleavedReader: %v", err)
			continue
		}
		for j, rdr := range []io.Reader{
			NewInterleavedReaderField(
				bytes.NewReader(encrypted.Bytes()[headerLen:]),
				QR_CODE_FIELD_256, eccLen,
			),
			rdr,
		} {
			hsh.Reset()
			read, err := io.Copy(hsh, rdr)
			if err != nil {
				t.Errorf("%d: error reading: %v", j, err)
			}
			if read != written {
				t.Errorf("%d: length mismatch: written %d, read %d", j, written, read)
			}

			readHash := hsh.Sum(nil)
			t.Logf("%d: read %d bytes, hash is %x", j, read, readHash)
			if !bytes.Equal(readHash, origHash) {
				t.Fatalf("%d: hash mismatch: written %x read %x", j, origHash, readHash)
			}
		}
	}
}
