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

import "fmt"

// ErrInsufficient is returned when there aren't enough bytes in the
// data slice to fulfill what the control byte specifies
var ErrInsufficient = fmt.Errorf("insufficient data supplied")

func decodeBlock(control byte, buf []byte) (results [4]uint32, n int) {
	blens := lookup[control]
	if len(buf) < int(blens[0]+blens[1]+blens[2]+blens[3]) {
		return results, 0
	}
	for ix, blen := range blens {
		for jx := uint8(0); jx < blen; jx++ {
			results[ix] <<= 8
			results[ix] |= uint32(buf[n])
			n++
		}
	}
	return results, n
}
