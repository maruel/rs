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
	"testing"
)

func TestGetPolyAlpha(t *testing.T) {
	for poly := 0x100; poly < 0x200; poly++ {
		for α := byte(0); α < 255; α++ {
			f := newField(poly, α)
			if f == nil {
				continue
			}
			poly2, α2 := f.GetPolyAlpha()
			if poly != poly2 || α != α2 {
				t.Errorf("got (%d, %d), awaited (%d, %d)", poly, α, poly2, α2)
			}
		}
	}
}
