// Copyright 2017 Nelz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svb

import "testing"

var offsets = []uint8{6, 4, 2, 0}

func TestLengths(t *testing.T) {
	for key, vals := range lengths {
		var total uint8
		for ix, offset := range offsets {
			expected := uint8(1 + ((key >> offset) & 0x03))
			total += expected
			if vals[ix] != expected {
				t.Errorf("%#x, %d: %d != %d\n", key, ix, expected, vals[ix])
			}
		}
		if total != vals[4] {
			t.Errorf("%#x, total: %d != %d\n", key, total, vals[4])
		}
	}
}
