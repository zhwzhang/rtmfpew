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

const SessionCloseAcknowledgementType = 0x4c

// SessionCloseAcknowledgement is sent in response to a SessionCloseRequest
type SessionCloseAcknowledgement struct{}

// Type returns SessionCloseAcknowledgement type opcode
func (chnk *SessionCloseAcknowledgement) Type() byte {
	return SessionCloseAcknowledgementType
}

func (chnk *SessionCloseAcknowledgement) Len() uint16 {
	return uint16(1)
}

func (chnk *SessionCloseAcknowledgement) WriteTo(buffer *bytes.Buffer) error {
	err := buffer.WriteByte(SessionCloseAcknowledgementType)
	if err != nil {
		return err
	}

	return binary.Write(buffer, binary.BigEndian, uint16(0))
}

func (chnk *SessionCloseAcknowledgement) ReadFrom(buffer *bytes.Buffer) error {
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	return err
}
