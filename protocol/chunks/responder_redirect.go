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
	"container/list"
	"encoding/binary"
	"github.com/rtmfpew/rtmfpew/protocol/connection"
	"github.com/rtmfpew/rtmfpew/protocol/vlu"
)

const ResponderRedirectChunkType = 0x71

// ResponderRedirectChunk is sent as response on InitiatorHello or ForwardedHello in startup mode
type ResponderRedirectChunk struct {
	TagEcho             []byte
	RedirectDestination []connection.PeerAddress
}

// Type returns ResponderRedirectChunk type opcode
func (chnk *ResponderRedirectChunk) Type() byte {
	return ResponderRedirectChunkType
}

func (chnk *ResponderRedirectChunk) Len() uint16 {
	TagEchoVlu := vlu.Vlu(len(chnk.TagEcho))
	destinationsLength := 0

	for _, destination := range chnk.RedirectDestination {
		destinationsLength += destination.Length()
	}

	return uint16(1 +
		TagEchoVlu.ByteLength() +
		len(chnk.TagEcho) +
		destinationsLength)
}

func (chnk *ResponderRedirectChunk) WriteTo(buffer *bytes.Buffer) error {

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

	for _, destination := range chnk.RedirectDestination {
		if err = destination.WriteTo(buffer); err != nil {
			return err
		}
	}

	return nil
}

func (chnk *ResponderRedirectChunk) ReadFrom(buffer *bytes.Buffer) error {
	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	tagLength := byte(0)
	tagLength, chnk.TagEcho, err = vlu.ReadVluBytesFrom(buffer)
	if err != nil {
		return err
	}

	destinationsLength := (int(length) - int(tagLength) - len(chnk.TagEcho))
	destinationsList := list.New()

	totalRead := 0
	for totalRead < destinationsLength {
		addr := &connection.PeerAddress{}
		if err = addr.ReadFrom(buffer); err != nil {
			return err
		}

		destinationsList.PushBack(addr)
		totalRead += (*addr).Length()
	}

	chnk.RedirectDestination = make([]connection.PeerAddress, destinationsList.Len())
	i := 0
	for dest := destinationsList.Front(); dest != nil; dest = dest.Next() {
		chnk.RedirectDestination[i] = *dest.Value.(*connection.PeerAddress)
		i++
	}

	return nil
}
