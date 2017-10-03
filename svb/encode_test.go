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

func TestEncodeBlock(t *testing.T) {
	tests := []struct {
		input   []uint32
		control byte
		data    []byte
	}{
		{
			[]uint32{1024, 12, 10, 1073741824},
			0x43, // 01 | 00 | 00 | 11
			[]byte{
				0x04, 0x00, // 1024
				0x0c,                   // 12
				0x0a,                   // 10
				0x40, 0x00, 0x00, 0x00, // 1,073,741,824
			},
		},
	}

	for _, test := range tests {
		d := make([]byte, 16)
		c, n := encodeBlock(d, test.input[0], test.input[1], test.input[2], test.input[3])

		if c != test.control {
			t.Errorf("control: %#x != %#x\n", c, test.control)
		}
		blens := lookup[c]
		size := blens[0] + blens[1] + blens[2] + blens[3]
		if n != int(size) {
			t.Errorf("size: %d != %d", n, size)
		}
	}
}
