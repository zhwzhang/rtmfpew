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

package session

import (
	"bytes"
	"container/list"
	"encoding/binary"

	"github.com/rtmfpew/amfy/vlu"
	"github.com/rtmfpew/rtmfpew/protocol/chunks"
	"io"
)

const (
	// ForbiddenMode should be ignored
	ForbiddenMode = 0
	// InitiatorMode used for session handshake
	InitiatorMode
	// ResponderMode used for communication
	ResponderMode
	// StartupMode used for session startup
	StartupMode
)

// Chunk generic interface
type Chunk interface {
	Type() byte
	Len() uint16
	ReadFrom(b *bytes.Buffer) error
	WriteTo(b *bytes.Buffer) error
}

// Packet contains chunks of data, and delivery options
type Packet struct {
	TimeCritical        bool
	TimeCriticalReserve bool

	TimestampPresent     bool
	TimestampEchoPresent bool

	Mode byte

	Timestamp     uint16
	TimestampEcho uint16

	HeaderLength uint32
	DataLength   uint32

	Chunks *list.List
}

// These Methods used in session facade

func (pckt *Packet) writeTo(buffer *bytes.Buffer) error {

	flags := byte(0)

	if pckt.TimeCritical {
		vlu.SetBit(&flags, 7)
	}

	if pckt.TimeCriticalReserve {
		vlu.SetBit(&flags, 6)
	}

	if pckt.TimestampPresent {
		vlu.SetBit(&flags, 3)
	}

	if pckt.TimestampEchoPresent {
		vlu.SetBit(&flags, 2)
	}

	flags = flags | pckt.Mode

	pckt.HeaderLength = 1
	if pckt.TimestampPresent {
		binary.Write(buffer, binary.BigEndian, pckt.Timestamp)
		pckt.HeaderLength++
	}

	if pckt.TimestampEchoPresent {
		binary.Write(buffer, binary.BigEndian, pckt.TimestampEcho)
		pckt.HeaderLength++
	}

	return nil
}

func (pckt *Packet) readFrom(buffer *bytes.Buffer) error {

	flags, err := buffer.ReadByte()
	if err != nil {
		return err
	}

	pckt.TimeCritical = vlu.BitIsSet(&flags, 7)
	pckt.TimeCriticalReserve = vlu.BitIsSet(&flags, 6)

	pckt.TimestampPresent = vlu.BitIsSet(&flags, 3)
	pckt.TimestampEchoPresent = vlu.BitIsSet(&flags, 2)

	pckt.Mode = flags & 0x03

	pckt.HeaderLength = 1
	if pckt.TimestampPresent {
		if err = binary.Read(buffer, binary.BigEndian, pckt.Timestamp); err != nil {
			return err
		}

		pckt.HeaderLength++
	}

	if pckt.TimestampEchoPresent {
		if err = binary.Read(buffer, binary.BigEndian, pckt.TimestampEcho); err != nil {
			return err
		}

		pckt.HeaderLength++
	}

	return nil
}

func (p *Packet) readChunks(buff *bytes.Buffer) (err error) {
	p.Chunks = list.New()
	var c Chunk
loop:
	for chunkType := byte(0); ; {
		if chunkType, err = buff.ReadByte(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		switch chunkType { // todo: reorder descending from most frequent
		case chunks.BufferProbeChunkType:
			c = &chunks.BufferProbeChunk{}
		case chunks.DataAcknowledgementBitmapChunkType:
			c = &chunks.DataAcknowledgementBitmapChunk{}
		case chunks.DataAcknowledgementRangesChunkType:
			c = &chunks.DataAcknowledgementRangesChunk{}
		case chunks.FlowExceptionReportChunkType:
			c = &chunks.FlowExceptionReportChunk{}
		case chunks.ForwardedHelloChunkType:
			c = &chunks.ForwardedHelloChunk{}
		case chunks.InitiatorHelloChunkType:
			c = &chunks.InitiatorHelloChunk{}
		case chunks.InitiatorInitialKeyingChunkType:
			c = &chunks.InitiatorInitialKeyingChunk{}
		case chunks.NextUserDataChunkType:
			c = &chunks.NextUserDataChunk{}
		case chunks.PingReplyChunkType:
			c = &chunks.PingReplyChunk{}
		case chunks.PingChunkType:
			c = &chunks.PingChunk{}
		case chunks.ResponderHelloChunkType:
			c = &chunks.ResponderHelloChunk{}
		case chunks.ResponderInitialKeyingChunkType:
			c = &chunks.ResponderInitialKeyingChunk{}
		case chunks.ResponderRedirectChunkType:
			c = &chunks.ResponderRedirectChunk{}
		case chunks.SessionCloseAcknowledgementType:
			c = &chunks.SessionCloseAcknowledgement{}
		case chunks.SessionCloseRequestChunkType:
			c = &chunks.SessionCloseRequestChunk{}
		case chunks.UserDataChunkType:
			c = &chunks.UserDataChunk{}
		case chunks.FragmentChunkType:
			c = &chunks.FragmentChunk{}
		default:
			break loop
		}
		if err = c.ReadFrom(buff); err != nil {
			return
		}
		p.Chunks.PushBack(c)
	}
	return
}

func (pckt *Packet) writeChunkTo(chnk Chunk, buffer *bytes.Buffer) error {
	lenBefore := buffer.Len()
	err := chnk.WriteTo(buffer)
	if err != nil {
		return err
	}

	pckt.DataLength += uint32(buffer.Len() - lenBefore)
	return nil
}

func (pckt *Packet) writePaddingTo(buffer *bytes.Buffer) error {
	padding := make([]byte, (pckt.DataLength+pckt.HeaderLength-1)%16)

	for i := uint32(0); i < (pckt.DataLength+pckt.HeaderLength-1)%16; i++ {
		padding[i] = 0xFF
	}

	return binary.Write(buffer, binary.BigEndian, padding)
}
