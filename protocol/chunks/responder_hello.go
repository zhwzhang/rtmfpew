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

	"github.com/rtmfpew/amfy/vlu"
)

const ResponderHelloChunkType = 0x70

// ResponderHelloChunk is sent as a response on InitiatorHelloChunk
type ResponderHelloChunk struct {
	TagEcho              []byte
	Cookie               []byte
	ResponderCertificate []byte
}

// Type returns ResponderHelloChunk type opcode
func (chnk *ResponderHelloChunk) Type() byte {
	return ResponderHelloChunkType
}

func (chnk *ResponderHelloChunk) Len() uint16 {
	TagEchoVlu := vlu.Vlu(len(chnk.TagEcho))
	CookieVlu := vlu.Vlu(len(chnk.Cookie))

	return uint16(1 +
		TagEchoVlu.ByteLength() +
		len(chnk.TagEcho) +
		CookieVlu.ByteLength() +
		len(chnk.Cookie) +
		len(chnk.ResponderCertificate))
}

func (chnk *ResponderHelloChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if err = vlu.WriteVluBytesTo(buffer, chnk.TagEcho); err != nil {
		return err
	}

	if err = vlu.WriteVluBytesTo(buffer, chnk.Cookie); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.ResponderCertificate); err != nil {
		return err
	}

	return nil
}

func (chnk *ResponderHelloChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}
	// Contents
	tagLength := byte(0)
	if tagLength, chnk.TagEcho, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	CookieLength := byte(0)
	if CookieLength, chnk.Cookie, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	certificateLength := (int(length) -
		int(tagLength) -
		len(chnk.TagEcho) -
		int(CookieLength) -
		len(chnk.Cookie))

	chnk.ResponderCertificate = make([]byte, certificateLength)
	num, err := buffer.Read(chnk.ResponderCertificate)

	if err != nil {
		return err
	}

	if num < int(certificateLength) {
		return err
	}

	return nil
}
