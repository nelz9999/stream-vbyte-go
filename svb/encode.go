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

var encodeOffsets = []uint8{24, 16, 8, 0}

func encodeBlock(control []byte, data []byte, num0, num1, num2, num3 uint32) int {
	i := 0
	for _, num := range []uint32{num0, num1, num2, num3} {
		length := 4
		if num < 256 { // 2**8
			length = 1
		} else if num < 65536 { // 2**16
			length = 2
		} else if num < 16777216 { // 2**24
			length = 3
		}

		control[0] = (control[0] << 2) + byte(length-1)
		for _, offset := range encodeOffsets[(4 - length):] {
			data[i] = byte(((num >> offset) & 0xff))
			i++
		}
	}

	return i
}
