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
	"github.com/rtmfpew/rtmfpew/protocol/io"
)

const ResponderHelloChunkType = 0x70

// ResponderHelloChunk is sent as a response on InitiatorHelloChunk
type ResponderHelloChunk struct {
	tagEcho              []byte
	cookie               []byte
	responderCertificate []byte
}

// Type returns ResponderHelloChunk type opcode
func (chnk *ResponderHelloChunk) Type() byte {
	return ResponderHelloChunkType
}

func (chnk *ResponderHelloChunk) Len() uint16 {
	tagEchoVlu := io.Vlu(len(chnk.tagEcho))
	cookieVlu := io.Vlu(len(chnk.cookie))
	
	return uint16(1 +
		tagEchoVlu.ByteLength() +
		len(chnk.tagEcho) +
		cookieVlu.ByteLength() +
		len(chnk.cookie) +
		len(chnk.responderCertificate))
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
	if err = io.WriteVluBytesTo(buffer, chnk.tagEcho); err != nil {
		return err
	}

	if err = io.WriteVluBytesTo(buffer, chnk.cookie); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.responderCertificate); err != nil {
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
	if tagLength, chnk.tagEcho, err = io.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	cookieLength := byte(0)
	if cookieLength, chnk.cookie, err = io.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	certificateLength := (int(length) -
		int(tagLength) -
		len(chnk.tagEcho) -
		int(cookieLength) -
		len(chnk.cookie))

	chnk.responderCertificate = make([]byte, certificateLength)
	num, err := buffer.Read(chnk.responderCertificate)

	if err != nil {
		return err
	}

	if num < int(certificateLength) {
		return err
	}

	return nil
}
