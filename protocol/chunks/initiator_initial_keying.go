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

const InitiatorInitialKeyingChunkType = 0x38

// InitiatorInitialKeyingChunk is sent as response on ResponderHelloChunk.
type InitiatorInitialKeyingChunk struct {
	InitiatorSessionID uint32
	// cookieLength vlu.Vlu
	CookieEcho []byte
	// certLength vlu.Vlu
	InitiatorCertificate []byte
	// initiatorComponentLength vlu.Vlu
	SessionKeyInitiatorComponent []byte

	Signature []byte
}

// Type returns InitiatorInitialKeyingChunk type opcode
func (chnk *InitiatorInitialKeyingChunk) Type() byte {
	return InitiatorInitialKeyingChunkType
}

func (chnk *InitiatorInitialKeyingChunk) Len() uint16 {
	cookieEchoVlu := vlu.Vlu(len(chnk.CookieEcho))
	initiatorCertVlu := vlu.Vlu(len(chnk.InitiatorCertificate))
	sessionKeyCompVlu := vlu.Vlu(len(chnk.SessionKeyInitiatorComponent))

	return uint16(1 +
		4 + // SessionID
		cookieEchoVlu.ByteLength() +
		len(chnk.CookieEcho) +
		initiatorCertVlu.ByteLength() +
		len(chnk.InitiatorCertificate) +
		sessionKeyCompVlu.ByteLength() +
		len(chnk.SessionKeyInitiatorComponent) +
		len(chnk.Signature))
}

func (chnk *InitiatorInitialKeyingChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if err = binary.Write(buffer, binary.BigEndian, chnk.InitiatorSessionID); err != nil {
		return err
	}

	if err = vlu.WriteVluBytesTo(buffer, chnk.CookieEcho); err != nil {
		return err
	}

	if err = vlu.WriteVluBytesTo(buffer, chnk.InitiatorCertificate); err != nil {
		return err
	}

	if err = vlu.WriteVluBytesTo(buffer, chnk.SessionKeyInitiatorComponent); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Signature); err != nil {
		return err
	}

	return nil
}

func (chnk *InitiatorInitialKeyingChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	if err = binary.Read(buffer, binary.BigEndian, &chnk.InitiatorSessionID); err != nil {
		return err
	}

	cookieEchoLength := byte(0)
	if cookieEchoLength, chnk.CookieEcho, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	initiatorCertificateLength := byte(0)
	if initiatorCertificateLength, chnk.InitiatorCertificate, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	sessionKeyInitiatorComponentLength := byte(0)
	if sessionKeyInitiatorComponentLength, chnk.SessionKeyInitiatorComponent, err = vlu.ReadVluBytesFrom(buffer); err != nil {
		return err
	}

	signatureLength := (int(length) -
		4 - // ID
		int(cookieEchoLength) -
		len(chnk.CookieEcho) -
		int(initiatorCertificateLength) -
		len(chnk.InitiatorCertificate) -
		int(sessionKeyInitiatorComponentLength) -
		len(chnk.SessionKeyInitiatorComponent))

	chnk.Signature = make([]byte, signatureLength)
	num, err := buffer.Read(chnk.Signature)

	if err != nil {
		return err
	}

	if num < signatureLength {
		return errors.New("Can't read initiator initial keying chunk signature")
	}

	return nil
}

// ChangeCookie creates HelloCookieChangeChunk to change cookie of InitiatorInitialKeyingChunk
func (keying *InitiatorInitialKeyingChunk) ChangeCookie(hello *ResponderHelloChunk) (freshCookie *HelloCookieChangeChunk, err error) {
	return nil, nil
}

// RespondKeying creates ResponderInitialKeyingChunk
func (keying *InitiatorInitialKeyingChunk) RespondKeying() (respondKeying *ResponderInitialKeyingChunk, err error) {
	return nil, nil
}
