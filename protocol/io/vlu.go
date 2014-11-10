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
	"errors"
	"math"
)

// Vlu uint32 wrapper
type Vlu uint32

func BitIsSet(value *byte, pos byte) bool {
	return (*value & (1 << (pos))) != 0
}

func SetBit(value *byte, pos byte) {
	*value = *value | (1 << (pos))
}

func get7BitValue(value uint32, pos byte) byte {
	return byte(((value) >> (7 * pos)) & uint32(math.Pow(2, 7)-1))
}

// ByteLength returns Vlu length for value in bytes
func (vlu *Vlu) ByteLength() int {

	if float64(*vlu) <= (math.Pow(2, 7) - 1) {
		return 1
	}

	if float64(*vlu) <= (math.Pow(2, 14) - 1) {
		return 2
	}

	if float64(*vlu) <= (math.Pow(2, 21) - 1) {
		return 3
	}

	if float64(*vlu) <= (math.Pow(2, 28) - 1) {
		return 4
	}

	return 0
}

func (vlu *Vlu) ReadFrom(buffer *bytes.Buffer) error {

	value := uint32(0)

	for total := 0; total <= 3; total++ {
		byt, _ := buffer.ReadByte()
		set := BitIsSet(&byt, 7)

		byt &^= 1 << 7 // Clear bit
		value <<= 7
		value |= uint32(byt)

		if !set {
			break
		} else if total == 3 {
			return errors.New("Vlu overflow")
		}
	}

	*vlu = Vlu(value)

	return nil
}

// ReadVluBytesFrom reads []byte array and returns it's Vlu size
func ReadVluBytesFrom(buffer *bytes.Buffer) (byte, []byte, error) {

	vlu := Vlu(0)
	if err := vlu.ReadFrom(buffer); err != nil {
		return 0, nil, err
	}

	data := make([]byte, vlu)

	num, err := buffer.Read(data)
	if Vlu(num) < vlu {
		return 0, nil, errors.New("Can't read VLU field properly")
	}

	if err != nil {
		return 0, nil, err
	}

	return byte(vlu.ByteLength()), data, nil
}

func (vlu *Vlu) WriteTo(buffer *bytes.Buffer) error {
	if uint32(*vlu) > uint32(math.Pow(2, 28) - 1) {
		return errors.New("Wrong VLU value")
	}

	if uint32(*vlu) == 0 {
		return buffer.WriteByte(byte(0))
	}

	for i := 3; i >= 0; i-- {
		valueChunk := get7BitValue(uint32(*vlu), byte(i))
		if valueChunk > 0 {
			if i != 0 {
				SetBit(&valueChunk, 7)
			}

			if err := buffer.WriteByte(valueChunk); err != nil {
				return err
			}
		}
	}

	return nil
}

func WriteVluBytesTo(buffer *bytes.Buffer, data []byte) error {
	vlu := Vlu(len(data))
	err := vlu.WriteTo(buffer)
	if err != nil {
		return err
	}

	_, err = buffer.Write(data)

	return err
}
