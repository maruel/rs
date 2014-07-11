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
	"code.google.com/p/rsc/gf256"
)

// The Galois Field for QR codes. See http://research.swtch.com/field for more
// information.
//
// x^8 + x^4 + x^3 + x^2 + 1
var QR_CODE_FIELD_256 = NewField(0x11D, 2)

// Field is a wrapper to gf256.Field so the type doesn't leak in.
type Field struct {
	f    *gf256.Field
	poly int
	α    byte
}

// NewField wraps gf256.NewField(). It is safe to use the premade
// QR_CODE_FIELD_256 all the time.
func NewField(poly int, α byte) *Field {
	return &Field{f: gf256.NewField(poly, int(α)), poly: poly, α: α}
}

func (f Field) GetPolyAlpha() (poly int, α byte) {
	if f.poly >= 0x100 && poly <= 0x200 {
		return f.poly, f.α
	}
	// guess
	α = f.f.Exp(1)
	for p := 0x100; p < 0x200; p++ {
		g := newField(p, α)
		if g == nil {
			continue
		}
		ok := true
		for i := 0; i < 255; i++ {
			if g.f.Exp(i) != f.f.Exp(i) {
				ok = false
				break
			}
		}
		if ok {
			return p, α
		}
	}
	return
}

func newField(poly int, α byte) (f *Field) {
	defer func() {
		recover()
	}()
	return NewField(poly, α)
}
