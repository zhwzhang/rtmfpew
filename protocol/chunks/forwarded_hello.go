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
	"github.com/rtmfpew/rtmfpew/protocol/connection"
	"github.com/rtmfpew/rtmfpew/protocol/vlu"
)

const ForwardedHelloChunkType = 0x0f

// ForwardedHelloChunk is a forwarded InitiatorHello
type ForwardedHelloChunk struct {
	Epd          []byte
	ReplyAddress connection.PeerAddress
	Tag          []byte
}

// Type returns ForwardedHelloChunk type opcode
func (chnk *ForwardedHelloChunk) Type() byte {
	return ForwardedHelloChunkType
}

func (chnk *ForwardedHelloChunk) Len() uint16 {
	v := vlu.Vlu(len(chnk.Epd))
	return uint16(1 + (&v).ByteLength() +
		len(chnk.Epd) +
		chnk.ReplyAddress.Length() +
		len(chnk.Tag))
}

func (chnk *ForwardedHelloChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
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

	if err = chnk.ReplyAddress.WriteTo(buffer); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Tag); err != nil {
		return err
	}

	return nil
}

func (chnk *ForwardedHelloChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	epdLength := byte(0)
	if epdLength, chnk.Epd, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	if err = chnk.ReplyAddress.ReadFrom(buffer); err != nil {
		return err
	}

	tagLength := int(length) -
		int(epdLength) -
		len(chnk.Epd) -
		chnk.ReplyAddress.Length()

	chnk.Tag = make([]byte, tagLength)
	num, err := buffer.Read(chnk.Tag)

	if err != nil {
		return err
	}

	if num < tagLength {
		return errors.New("Can't read ForwardedHello chunk tag")
	}

	return nil
}
