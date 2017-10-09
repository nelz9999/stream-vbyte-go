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

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestUint32s(t *testing.T) {
	tests := []struct {
		control byte
		data    []byte
		results []uint32
		isErr   bool
	}{
		// Figure 3 from https://arxiv.org/pdf/1709.08990.pdf
		{
			0x43, // 01 | 00 | 00 | 11
			[]byte{
				0x04, 0x00, // 1024
				0x0c,                   // 12
				0x0a,                   // 10
				0x40, 0x00, 0x00, 0x00, // 1,073,741,824
			},
			[]uint32{
				1024,
				12,
				10,
				1073741824,
			},
			false,
		},
		// Figure 3 from https://arxiv.org/pdf/1709.08990.pdf
		{
			0x01, // 00 | 00 | 00 | 01
			[]byte{
				0x01,       // 1
				0x02,       // 2
				0x03,       // 3
				0x04, 0x00, // 1024
			},
			[]uint32{
				1,
				2,
				3,
				1024,
			},
			false,
		},
		// This is an error case, expecting ErrInsufficient
		{
			0x00, // 00 | 00 | 00 | 00
			[]byte{0x00, 0x00, 0x00}, // insufficient
			[]uint32{0, 0, 0, 0},
			true,
		},
	}

	for _, test := range tests {
		// raw buffer form //
		// - - - - - - - - //
		r, n := Uint32s(test.control, test.data)

		if test.isErr {
			if n != 0 {
				t.Errorf("%#x: %d != 0\n", test.control, n)
			}
		} else {
			blens := lookup[test.control]
			size := int(blens[0] + blens[1] + blens[2] + blens[3])
			if n != size {
				t.Errorf("%#x: %d != %d\n", test.control, n, size)
			}
		}

		for ix, expected := range test.results {
			if r[ix] != expected {
				t.Errorf("%#x: %d != %d\n", test.control, r[ix], expected)
			}
		}

		// via io.ByteReader //
		// - - - - - - - - - //
		var err error
		r, err = ReadUint32s(test.control, bytes.NewBuffer(test.data))
		if err != nil && !test.isErr {
			t.Errorf("unexpected: %v\n", err)
		}

		for ix, expected := range test.results {
			if r[ix] != expected {
				t.Errorf("%#x: %d != %d\n", test.control, r[ix], expected)
			}
		}
	}
}

func TestGetU32BlockPanicForData(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic received")
		}
	}()
	GetU32Block(0xff, []byte{0x00, 0x00, 0x00}, false)
}

func TestGetU32Block(t *testing.T) {
	tests := []struct {
		ctrl byte
		data []byte
		quad []uint32
		size int
	}{
		{ // Smallest possible block
			0x00,
			[]byte{0x00, 0x00, 0x00, 0x00},
			[]uint32{0, 0, 0, 0},
			4,
		},
		{ // Smallest all-4-byte representation
			0xff,
			[]byte{
				0x01, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00,
			},
			[]uint32{(1 << 24), 2 * (1 << 24), 3 * (1 << 24), 4 * (1 << 24)},
			16,
		},
		{ // From whitepapaer: 1024, 12, 10, 1073741824
			0x43,
			[]byte{
				0x04, 0x00, // 1024
				0x0c,                   // 12
				0x0a,                   // 10
				0x40, 0x00, 0x00, 0x00, // 1073741824
			},
			[]uint32{1024, 1036, 1046, 1073742870},
			8,
		},
		{ // From whitepapaer: 1, 2, 3, 1024
			0x01,
			[]byte{
				0x01,       // 1
				0x02,       // 2
				0x03,       // 3
				0x04, 0x00, // 1024
			},
			[]uint32{1, 3, 6, 1030},
			5,
		},
	}

	for _, test := range tests {
		quad, size := GetU32Block(test.ctrl, test.data, true)
		if size != test.size {
			t.Errorf("mismatch size: %d != %d\n", size, test.size)
		}
		for ix := 0; ix < 4; ix++ {
			if quad[ix] != test.quad[ix] {
				t.Errorf("mismatch via %d: %x != %x\n", ix, quad[ix], test.quad[ix])
			}
		}
	}
}

func TestU32BlockRoundtrip(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]byte, 16)
	quad := []uint32{0, 0, 0, 0}
	for ix := 0; ix < 100; ix++ {
		diff := (ix%2 == 0)

		ctrl, size := PutU32Block(data, quad, diff)
		q, s := GetU32Block(ctrl, data[:size], diff)

		if s != size {
			t.Errorf("mismatch size: %d != %d for %v\n", s, size, quad)
		}

		// Test this block, but also setup the random data for
		// the next iteration
		for jx := range quad {
			if q[jx] != quad[jx] {
				t.Errorf("mismatch: %v != %v\n", q, quad)
			}

			// Generate random data for the next time round
			blen := uint(1 + r.Intn(4))
			top := uint32(1) << (8 * blen)
			low := (uint32(1) << (8 * (blen - 1))) - 1
			quad[jx] = low + uint32(r.Intn(int(top-low)))
			// t.Logf("%d: %d\n", jx, quad[jx])
		}
		// t.Logf("size: %d\n", size)
	}
}
