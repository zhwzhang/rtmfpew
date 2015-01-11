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

	"github.com/rtmfpew/amfy/vlu"
)

const ResponderInitialKeyingChunkType = 0x78

// ResponderInitialKeyingChunk is sent as response on InitiatorInitialKeying in startup mode
type ResponderInitialKeyingChunk struct {
	ResponderSessionID uint32
	// responderComponentLength vlu.Vlu
	SessionKeyResponderComponent []byte

	Signature []byte
}

// Type returns ResponderInitialKeyingChunk type opcode
func (chnk *ResponderInitialKeyingChunk) Type() byte {
	return ResponderInitialKeyingChunkType
}

func (chnk *ResponderInitialKeyingChunk) Len() uint16 {
	sessionKeyVlu := vlu.Vlu(len(chnk.SessionKeyResponderComponent))

	return uint16(1 +
		4 + // ID
		sessionKeyVlu.ByteLength() +
		len(chnk.SessionKeyResponderComponent) +
		len(chnk.Signature))
}

func (chnk *ResponderInitialKeyingChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if err = binary.Write(buffer, binary.BigEndian, chnk.ResponderSessionID); err != nil {
		return err
	}

	if err = vlu.WriteVluBytesTo(buffer, chnk.SessionKeyResponderComponent); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Signature); err != nil {
		return err
	}

	return nil
}

func (chnk *ResponderInitialKeyingChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	if err = binary.Read(buffer, binary.BigEndian, &chnk.ResponderSessionID); err != nil {
		return err
	}

	sessionKeyResponderComponentLength := byte(0)
	if sessionKeyResponderComponentLength, chnk.SessionKeyResponderComponent, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	signatureLength := (int(length) -
		4 - // ID
		int(sessionKeyResponderComponentLength) -
		len(chnk.SessionKeyResponderComponent))

	chnk.Signature = make([]byte, signatureLength)
	num, err := buffer.Read(chnk.Signature)

	if err != nil {
		return err
	}

	if num < signatureLength {
		return errors.New("Can't read responder initial keying chunk signature")
	}

	return nil
}
