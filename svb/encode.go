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

func encodeBlock(buf []byte, num0, num1, num2, num3 uint32) (control byte, n int) {
	for _, num := range []uint32{num0, num1, num2, num3} {
		control <<= 2
		blen := byteLength(num)
		control |= byte(blen - 1)
		for _, offset := range offsets[(4 - blen):] {
			buf[n] = byte((num >> offset) & 0xff)
			n++
		}
	}
	return control, n
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
