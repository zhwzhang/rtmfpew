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
	"github.com/rtmfpew/rtmfpew/protocol/io"
)

const HelloCookieChangeChunkType = 0x79

// HelloCookieChangeChunk is sent to change cookie of InitiatorInitialKeying in startup mod
type HelloCookieChangeChunk struct {
	// OldCookieLen io.Vlu
	OldCookie []byte
	NewCookie []byte
}

// Type returns HelloCookieChangeChunk type opcode
func (chnk *HelloCookieChangeChunk) Type() byte {
	return HelloCookieChangeChunkType
}

func (chnk *HelloCookieChangeChunk) Len() uint16 {
	OldCookieVlu := io.Vlu(len(chnk.OldCookie))

	return uint16(1 +
		(&OldCookieVlu).ByteLength() +
		len(chnk.OldCookie) +
		len(chnk.NewCookie))
}

func (chnk *HelloCookieChangeChunk) WriteTo(buffer *bytes.Buffer) error {
	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	if err = io.WriteVluBytesTo(buffer, chnk.OldCookie); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.NewCookie); err != nil {
		return err
	}

	return nil
}

func (chnk *HelloCookieChangeChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	oldCookieLength := byte(0)
	oldCookieLength, chnk.OldCookie, err = io.ReadVluBytesFrom(buffer)
	if err != nil {
		return err
	}

	newCookieLength := (int(length) - int(oldCookieLength) - len(chnk.OldCookie))
	chnk.NewCookie = make([]byte, newCookieLength)

	num, err := buffer.Read(chnk.NewCookie)

	if err != nil {
		return err
	}

	if num < newCookieLength {
		return errors.New("Can't read hello cookie change chunk")
	}

	return nil
}
