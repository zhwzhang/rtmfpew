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
)

const SessionCloseRequestChunkType = 0x0c

// SessionCloseRequestChunk is sent to terminate a session
type SessionCloseRequestChunk struct{}

// Type returns SessionCloseRequestChunk type opcode
func (chnk *SessionCloseRequestChunk) Type() byte {
	return SessionCloseRequestChunkType
}

func (chnk *SessionCloseRequestChunk) Len() uint16 {
	return uint16(1)
}

func (chnk *SessionCloseRequestChunk) WriteTo(buffer *bytes.Buffer) error {
	err := buffer.WriteByte(SessionCloseRequestChunkType)
	if err != nil {
		return err
	}

	return binary.Write(buffer, binary.BigEndian, uint16(0))
}

func (chnk *SessionCloseRequestChunk) ReadFrom(buffer *bytes.Buffer) error {
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	return err
}
