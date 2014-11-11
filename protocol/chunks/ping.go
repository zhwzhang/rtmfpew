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

package chunks

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const PingChunkType = 0x01

// PingChunk is a simple ping Message
type PingChunk struct {
	Message []byte
}

// Type returns PingChunk type opcode
func (chnk *PingChunk) Type() byte {
	return PingChunkType
}

func (chnk *PingChunk) Len() uint16 {
	return uint16(1 + len(chnk.Message))
}

func (chnk *PingChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if _, err = buffer.Write(chnk.Message); err != nil {
		return err
	}

	return nil
}

func (chnk *PingChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	chnk.Message = make([]byte, length)
	num, err := buffer.Read(chnk.Message)

	if err != nil {
		return err
	}

	if num < int(length) {
		return errors.New("Can't read ping chunk Message")
	}

	return nil
}
