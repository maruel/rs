/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Original source:
// https://code.google.com/p/zxing/source/browse/trunk/core/test/src/com/google/zxing/common/reedsolomon/ReedSolomonDecoderQRCodeTestCase.java
//
// Copyright 2008 ZXing authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @author Sean Owen

package rs

import (
	"math/rand"
	"testing"

	"github.com/willf/bitset"
)

// See ISO 18004, Appendix I, from which this example is taken.
var QRCodeTestData = []byte{0x10, 0x20, 0x0C, 0x56, 0x61, 0x80, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11}
var QRCodeTestECC = []byte{0xA5, 0x24, 0xD4, 0xC1, 0xED, 0x36, 0xC7, 0x87, 0x2C, 0x55}
var QRCodeCorrectable = len(QRCodeTestECC) / 2

// 128 pseudo random bytes.
var Rand128 = []byte{
	0x1b, 0x01, 0xf2, 0xd6, 0x3b, 0x31, 0x2a, 0x87, 0x6c, 0x93, 0x96, 0x11, 0x13, 0x8f, 0x20, 0x06,
	0x6d, 0x70, 0x26, 0x6e, 0xc2, 0x76, 0xf0, 0xef, 0x9a, 0xda, 0x4d, 0xbe, 0x71, 0xb7, 0x6c, 0xbb,
	0x4d, 0x0b, 0xc2, 0x27, 0x2a, 0x5b, 0xbf, 0x79, 0xb9, 0xc0, 0xc8, 0xef, 0x24, 0x7c, 0x9d, 0xb3,
	0x08, 0x6d, 0xe7, 0x09, 0x54, 0x8d, 0x13, 0xda, 0x40, 0x21, 0x27, 0x3a, 0x39, 0x44, 0x34, 0x2d,
	0xab, 0xa1, 0x6a, 0x40, 0x79, 0xcd, 0x7c, 0x56, 0x31, 0xc2, 0x0f, 0x41, 0x02, 0x13, 0x17, 0x3e,
	0xcc, 0xef, 0x4c, 0x4e, 0xb9, 0xc2, 0x04, 0x7a, 0x58, 0x2a, 0x27, 0xa4, 0x92, 0x09, 0xdf, 0x20,
	0x22, 0x5c, 0xfb, 0x45, 0x47, 0x2e, 0xaa, 0x88, 0x51, 0xa1, 0xa1, 0x3d, 0xc1, 0x34, 0xbc, 0x34,
	0xe3, 0x53, 0x56, 0xb7, 0xd6, 0x43, 0x92, 0xf9, 0x47, 0xe4, 0xa9, 0xa1, 0x94, 0xad, 0x1a, 0x7f,
}

func makecopy(a []byte) []byte {
	b := make([]byte, len(a))
	copy(b, a)
	return b
}

// Returns the test data and the ECC codes as one slice.
func getcomplete() []byte {
	data := make([]byte, len(QRCodeTestData)+len(QRCodeTestECC))
	copy(data, QRCodeTestData)
	copy(data[len(QRCodeTestData):], QRCodeTestECC)
	return data
}

func TestNoError(t *testing.T) {
	data := makecopy(QRCodeTestData)
	ecc := makecopy(QRCodeTestECC)
	checkQR(t, data, ecc, 0)
}

func TestOneErrorData(t *testing.T) {
	for i := 0; i < len(QRCodeTestData); i++ {
		data := makecopy(QRCodeTestData)
		ecc := makecopy(QRCodeTestECC)
		data[i] = data[i] + byte(i+1)
		checkQR(t, data, ecc, 1)
	}
}

func TestOneErrorECC(t *testing.T) {
	for i := 0; i < len(QRCodeTestECC); i++ {
		data := makecopy(QRCodeTestData)
		ecc := makecopy(QRCodeTestECC)
		ecc[i] = ecc[i] + byte(i+1)
		checkQR(t, data, ecc, 1)
	}
}

func TestMaxErrors(t *testing.T) {
	complete := getcomplete()
	corrupt(complete, QRCodeCorrectable)
	checkQR(t, complete[:len(QRCodeTestData)], complete[len(QRCodeTestData):], QRCodeCorrectable)
}

func TestTooManyErrors(t *testing.T) {
	complete := getcomplete()
	corrupt(complete, QRCodeCorrectable+1)
	d := NewDecoder(QR_CODE_FIELD_256)
	if nb, err := d.Decode(complete[:len(QRCodeTestData)], complete[len(QRCodeTestData):]); err == nil {
		t.Fatal("Recovered unrecoverable error!?!")
	} else if nb != 0 {
		t.Fatalf("Err != %d", nb)
	}
}

func checkQR(t *testing.T, data, ecc []byte, nbErrors int) {
	goldenData := QRCodeTestData
	goldenEcc := QRCodeTestECC
	d := NewDecoder(QR_CODE_FIELD_256)
	errorsFound, err := d.Decode(data, ecc)
	if err != nil {
		t.Fatalf("Got error: %s", err)
	}
	if nbErrors != errorsFound {
		t.Fatalf("Expected %d errors, got %d", nbErrors, errorsFound)
	}
	compare(t, data, goldenData, "Data differs")
	compare(t, ecc, goldenEcc, "ECC differs")
}

// https://code.google.com/p/zxing/source/browse/trunk/core/test/src/com/google/zxing/common/reedsolomon/AbstractReedSolomonTestCase.java
func corrupt(received []byte, howMany int) {
	corrupted := bitset.New(uint(len(received)))
	for j := 0; j < howMany; j++ {
		location := uint(rand.Int31n(int32(len(received))))
		if corrupted.Test(location) {
			j--
		} else {
			corrupted.Set(location)
			received[location] = (received[location] + 1 + byte(rand.Int31n(255))) & 0xFF
		}
	}
}

func BenchmarkDecode16_10(b *testing.B) {
	b.StopTimer()
	d := NewDecoder(QR_CODE_FIELD_256)
	data := makecopy(QRCodeTestData)
	ecc := makecopy(QRCodeTestECC)
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n, err := d.Decode(data, ecc); err != nil || n != 0 {
			b.Fail()
		}
	}
}

func BenchmarkDecode16_10With1Error(b *testing.B) {
	b.StopTimer()
	d := NewDecoder(QR_CODE_FIELD_256)
	data := makecopy(QRCodeTestData)
	data[1] = data[1] + 1
	ecc := makecopy(QRCodeTestECC)
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n, err := d.Decode(makecopy(data), ecc); err != nil || n != 1 {
			b.Fail()
		}
	}
}

func BenchmarkDecode128_16(b *testing.B) {
	b.StopTimer()
	data := makecopy(Rand128)
	ecc := make([]byte, 16)
	e := NewEncoder(QR_CODE_FIELD_256, len(ecc))
	e.Encode(data, ecc)
	d := NewDecoder(QR_CODE_FIELD_256)
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n, err := d.Decode(data, ecc); err != nil || n != 0 {
			b.Fail()
		}
	}
}

func BenchmarkDecode128_16With1Error(b *testing.B) {
	b.StopTimer()
	data := makecopy(Rand128)
	ecc := make([]byte, 16)
	e := NewEncoder(QR_CODE_FIELD_256, len(ecc))
	e.Encode(data, ecc)
	d := NewDecoder(QR_CODE_FIELD_256)
	data[1] = data[1] + 1
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n, err := d.Decode(makecopy(data), ecc); err != nil || n != 1 {
			b.Fail()
		}
	}
}

func BenchmarkDecode128_16With2Error(b *testing.B) {
	b.StopTimer()
	data := makecopy(Rand128)
	ecc := make([]byte, 16)
	e := NewEncoder(QR_CODE_FIELD_256, len(ecc))
	e.Encode(data, ecc)
	d := NewDecoder(QR_CODE_FIELD_256)
	data[1] = data[1] + 1
	ecc[1] = ecc[1] + 1
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n, err := d.Decode(makecopy(data), makecopy(ecc)); err != nil || n != 2 {
			b.Fail()
		}
	}
}
