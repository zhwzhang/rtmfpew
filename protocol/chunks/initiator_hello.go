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

const InitiatorHelloChunkType = 0x30

// InitiatorHelloChunk initiates RTMFP handshake in startup mode
type InitiatorHelloChunk struct {
	Epd []byte
	Tag []byte
}

// Type returns InitiatorHelloChunk type opcode
func (chnk *InitiatorHelloChunk) Type() byte {
	return InitiatorHelloChunkType
}

func (chnk *InitiatorHelloChunk) Len() uint16 {
	v := vlu.Vlu(len(chnk.Epd))
	return uint16(1 +
		(&v).ByteLength() +
		len(chnk.Epd) +
		len(chnk.Tag))
}

func (chnk *InitiatorHelloChunk) WriteTo(buffer *bytes.Buffer) error {

	err := buffer.WriteByte(chnk.Type())
	// Chunk header
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if err = vlu.WriteVluBytesTo(buffer, chnk.Epd); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Tag); err != nil {
		return err
	}

	return nil
}

func (chnk *InitiatorHelloChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	EpdLength := byte(0)

	// Contents
	if EpdLength, chnk.Epd, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	TagLength := int(length) - int(EpdLength) - len(chnk.Epd)
	chnk.Tag = make([]byte, TagLength)

	num, err := buffer.Read(chnk.Tag)

	if err != nil {
		return err
	}

	if num < TagLength {
		return errors.New("Can't read initiator hello chunk tag")
	}

	return nil
}
