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
	"github.com/willf/bitset"
	"math/rand"
	"testing"
)

// See ISO 18004, Appendix I, from which this example is taken.
var QR_CODE_TEST_DATA = []byte{0x10, 0x20, 0x0C, 0x56, 0x61, 0x80, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11, 0xEC, 0x11}
var QR_CODE_TEST_ECC = []byte{0xA5, 0x24, 0xD4, 0xC1, 0xED, 0x36, 0xC7, 0x87, 0x2C, 0x55}
var QR_CODE_CORRECTABLE = len(QR_CODE_TEST_ECC) / 2

func makecopy(a []byte) []byte {
	b := make([]byte, len(a))
	copy(b, a)
	return b
}

// Returns the test data and the ECC codes as one slice.
func getcomplete() []byte {
	data := make([]byte, len(QR_CODE_TEST_DATA)+len(QR_CODE_TEST_ECC))
	copy(data, QR_CODE_TEST_DATA)
	copy(data[len(QR_CODE_TEST_DATA):], QR_CODE_TEST_ECC)
	return data
}

func TestNoError(t *testing.T) {
	data := makecopy(QR_CODE_TEST_DATA)
	ecc := makecopy(QR_CODE_TEST_ECC)
	checkQR(t, data, ecc, 0)
}

func TestOneErrorData(t *testing.T) {
	for i := 0; i < len(QR_CODE_TEST_DATA); i++ {
		data := makecopy(QR_CODE_TEST_DATA)
		ecc := makecopy(QR_CODE_TEST_ECC)
		data[i] = data[i] + byte(i+1)
		checkQR(t, data, ecc, 1)
	}
}

func TestOneErrorECC(t *testing.T) {
	for i := 0; i < len(QR_CODE_TEST_ECC); i++ {
		data := makecopy(QR_CODE_TEST_DATA)
		ecc := makecopy(QR_CODE_TEST_ECC)
		ecc[i] = ecc[i] + byte(i+1)
		checkQR(t, data, ecc, 1)
	}
}

func TestMaxErrors(t *testing.T) {
	complete := getcomplete()
	corrupt(complete, QR_CODE_CORRECTABLE)
	checkQR(t, complete[:len(QR_CODE_TEST_DATA)], complete[len(QR_CODE_TEST_DATA):], QR_CODE_CORRECTABLE)
}

func TestTooManyErrors(t *testing.T) {
	complete := getcomplete()
	corrupt(complete, QR_CODE_CORRECTABLE+1)
	d:= NewDecoder(QR_CODE_FIELD_256)
	if nb, err := d.Decode(complete[:len(QR_CODE_TEST_DATA)], complete[len(QR_CODE_TEST_DATA):]); err == nil {
		t.Fatal("Recovered unrecoverable error!?!")
	} else if nb != 0 {
		t.Fatal("Err != %d", nb)
	}
}

func checkQR(t *testing.T, data, ecc []byte, nbErrors int) {
	goldenData := QR_CODE_TEST_DATA
	goldenEcc := QR_CODE_TEST_ECC
	d:= NewDecoder(QR_CODE_FIELD_256)
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
	data := makecopy(QR_CODE_TEST_DATA)
	ecc := makecopy(QR_CODE_TEST_ECC)
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
	data := makecopy(QR_CODE_TEST_DATA)
	data[1] = data[1] + 1
	ecc := makecopy(QR_CODE_TEST_ECC)
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
	data := make([]byte, 128)
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
	data := make([]byte, 128)
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
	data := make([]byte, 128)
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
