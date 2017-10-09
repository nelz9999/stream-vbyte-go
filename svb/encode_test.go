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

func TestPutUint32s(t *testing.T) {
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
		c, n := PutUint32s(d, test.input[0], test.input[1], test.input[2], test.input[3])

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

func TestPutUint32sPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic received")
		}
	}()
	PutUint32s([]byte{0x00}, 0, 1, 2, 3)
}

func TestPutU32BlockPanicForData(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic received")
		}
	}()
	PutU32Block([]byte{}, []uint32{0, 0, 0, 0}, false)
}

func TestPutU32BlockPanicForQuad(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic received")
		}
	}()
	data := make([]byte, 16)
	PutU32Block(data, []uint32{0}, false)
}

func TestPutU32Block(t *testing.T) {
	tests := []struct {
		quad []uint32
		ctrl byte
		size int
	}{
		{ // Smallest possible encoded
			[]uint32{0, 0, 0, 0},
			0x00,
			4,
		},
		{ // Smallest all-4-byte representation
			[]uint32{(1 << 24), 2 * (1 << 24), 3 * (1 << 24), 4 * (1 << 24)},
			0xff,
			16,
		},
		{ // From whitepapaer (after diff coding): 1024, 12, 10, 1073741824
			[]uint32{1024, 1036, 1046, 1073742870},
			0x43,
			8,
		},
		{ // From whitepapaer (after diff coding): 1, 2, 3, 1024
			[]uint32{1, 3, 6, 1030},
			0x01,
			5,
		},
	}

	for _, test := range tests {
		data := make([]byte, 16)
		ctrl, size := PutU32Block(data, test.quad, true)
		if ctrl != test.ctrl {
			t.Errorf("ctrl mismatch: %x != %x\n", ctrl, test.ctrl)
		}
		if size != test.size {
			t.Errorf("size mismatch: %d != %d\n", size, test.size)
		}
		// t.Logf("% x\n", data[:size])
	}
}
