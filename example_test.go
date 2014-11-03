/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

package rs

import (
	"fmt"
)

func ExampleNewEncode() {
	data := []byte("hello, world")
	fmt.Printf("Original data: %s\n", data)
	ecc := make([]byte, 2)
	e := NewEncoder(QRCodeField256, len(ecc))
	e.Encode(data, ecc)
	fmt.Printf("ECC bytes: %v\n", ecc)
	// Output:
	// Original data: hello, world
	// ECC bytes: [171 167]
}

func ExampleNewDecode() {
	data := []byte("hello, wXrld")
	ecc := []byte{171, 167}
	fmt.Printf("Corrupted data: %s\n", data)
	d := NewDecoder(QRCodeField256)
	if nb, err := d.Decode(data, ecc); err != nil || nb != 1 {
		fmt.Printf("Expected 1 fix, for %d. Error: %s\n", nb, err)
	}
	fmt.Printf("Fixed data: %s\n", data)
	// Output:
	// Corrupted data: hello, wXrld
	// Fixed data: hello, world
}
