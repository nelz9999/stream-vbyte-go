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

var offsets = []uint8{24, 16, 8, 0}

// PutUint32s encodes a quad of uint32 into the data buffer, returning
// the control byte that signifies the encoded byte lengths, and the length
// of how many bytes got written to the data buffer.
// If the buffer is too small, PutUStreamVByte will panic
func PutUint32s(data []byte, num0, num1, num2, num3 uint32) (ctrl byte, n int) {
	for _, num := range []uint32{num0, num1, num2, num3} {
		ctrl <<= 2
		blen := byteLength(num)
		ctrl |= byte(blen - 1)
		for _, offset := range offsets[(4 - blen):] {
			data[n] = byte((num >> offset) & 0xff)
			n++
		}
	}
	return ctrl, n
}

func byteLength(n uint32) uint8 {
	if n < 256 {
		return 1
	}
	if n < 65536 {
		return 2
	}
	if n < 16777216 {
		return 3
	}
	return 4
}

// PutU32Block encodes a single quad of uint32 values. (This function is the
// write-side parallel to GetU32Block. These operations are optimized for
// read-side speed, but the write-side is still pretty quick.)
//
// The data parameter is the buffer where the encoded values are written, and
// may need up to 16 bytes available. The quad parameter is the buffer of
// integers that are to be encoded, and needs to have 4 values available.
// The diff value signifies that you want to use "differential coding" (for
// more efficient storage of the values), but this requires that the values
// must be in ascending sorted order.
//
// The ctrl byte returned is a required hint for decoding, and is expected to
// be stored in parallel to the data buffer. And the return value n represents
// the number of bytes used in the data buffer, as per the algorithm.
//
// Panics will be thrown if there are too few bytes available in the data
// buffer, or too few values in the quad buffer.
func PutU32Block(data []byte, quad []uint32, diff bool) (ctrl byte, n int) {
	var prev uint32
	for i := uint(0); i < 4; i++ {
		num := quad[i]
		if diff {
			num = num - prev
			prev += num
		}
		blen := byteLength(num)
		ctrl |= ((blen - 1) << (6 - 2*i))
		if blen == 4 {
			data[n] = byte((num >> 24) & 0xff)
			n++
		}
		if blen >= 3 {
			data[n] = byte((num >> 16) & 0xff)
			n++
		}
		if blen >= 2 {
			data[n] = byte((num >> 8) & 0xff)
			n++
		}
		data[n] = byte(num & 0xff)
		n++
	}
	return ctrl, n
}
