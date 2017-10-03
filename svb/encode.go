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
