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

import "io"

// Uint32s decodes a quad of uint32 from the data buffer, returning
// the four uint32s and the number of bytes consumed from the buffer.
// If there aren't enough bytes in the data buffer to match what is required
// by the ctrl byte, the quad is returned as zeros with n = 0.
func Uint32s(ctrl byte, data []byte) (nums [4]uint32, n int) {
	blens := lookup[ctrl]
	if len(data) < int(blens[0]+blens[1]+blens[2]+blens[3]) {
		return nums, 0
	}
	for ix, blen := range blens {
		for jx := uint8(0); jx < blen; jx++ {
			nums[ix] <<= 8
			nums[ix] |= uint32(data[n])
			n++
		}
	}
	return nums, n
}

// ReadUint32s reads a quad of uint32 from d, using the information encoded
// in the ctrl byte.
func ReadUint32s(ctrl byte, d io.ByteReader) (nums [4]uint32, err error) {
	blens := lookup[ctrl]
	var n int
	for ix, blen := range blens {
		for jx := uint8(0); jx < blen; jx++ {
			b, err := d.ReadByte()
			if err != nil {
				return nums, err
			}
			nums[ix] <<= 8
			nums[ix] |= uint32(b)
			n++
		}
	}
	return nums, nil
}
