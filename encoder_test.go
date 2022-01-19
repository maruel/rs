/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Original source:
// https://github.com/zxing/zxing/tree/master/core/src/test/java/com/google/zxing/common/reedsolomon
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
	"bytes"
	"testing"
)

// Tests example given in ISO 18004, Annex I.
func TestISO18004Example(t *testing.T) {
	actualEcc := make([]byte, len(QRCodeTestECC))
	NewEncoder(QRCodeField256, len(QRCodeTestECC)).Encode(QRCodeTestData, actualEcc)
	compare(t, QRCodeTestECC, actualEcc, "ECC differs")
}

func compare(t *testing.T, a []byte, b []byte, msg string) {
	if !bytes.Equal(a, b) {
		t.Fatalf("%s: %q != %q", msg, a, b)
	}
}

// Sample QR Code.
func BenchmarkEncode16_10(b *testing.B) {
	b.StopTimer()
	data := makecopy(Rand128[:16])
	ecc := [10]byte{}
	e := NewEncoder(QRCodeField256, len(ecc))
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(data[:], ecc[:])
	}
}

// 12.5% ECC size.
func BenchmarkEncode128_16(b *testing.B) {
	b.StopTimer()
	data := makecopy(Rand128)
	ecc := [16]byte{}
	e := NewEncoder(QRCodeField256, len(ecc))
	b.SetBytes(int64(len(data) * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		e.Encode(data[:], ecc[:])
	}
}
