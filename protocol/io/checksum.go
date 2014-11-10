//
// Copyright 2014 RTMFPew
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package io

import (
	"bytes"
	"encoding/binary"
)

// Checksum calculates IPv4 Header checksum from rfc1071
func Checksum(b []byte) uint16 {
	a := uint16(0)
	buff := bytes.NewBuffer(b)
	err := binary.Read(buff, binary.BigEndian, &a)
	if err != nil {
		return 0
	}

	acc := uint32(0)

	for i := 0; i < len(b)/2; i++ {
		err := binary.Read(buff, binary.BigEndian, &a)
		if err != nil {
			return 0
		}

		acc += uint32(a)
	}

	if len(b)%2 > 0 {
		t, err := buff.ReadByte()
		if err != nil {
			return 0
		}

		acc += uint32(t)
	}

	for acc>>16 != 0 {
		acc = (acc & 0xffff) + (acc >> 16)
	}

	return ^uint16(acc)
}
