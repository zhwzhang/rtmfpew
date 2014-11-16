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
	"github.com/rtmfpew/amfy/vlu"
)

const DataAcknowledgementRangesChunkType = 0x51

// DataAcknowledgementRange indicates what UserData fragments had been recieved.
// Is a part of DataAcknowledgementRangeChunk
type DataAcknowledgementRange struct {
	HolesMinusOne    vlu.Vlu
	ReceivedMinusOne vlu.Vlu
}

// DataAcknowledgementRangesChunk is sent to indicate UserData fragments have been recieved for one flow.
type DataAcknowledgementRangesChunk struct {
	FlowID                vlu.Vlu
	BufferBlocksAvailable vlu.Vlu
	CumulativeAck         vlu.Vlu

	Ranges []DataAcknowledgementRange
}

// Type returns DataAcknowledgementRangesChunk type opcode
func (chnk *DataAcknowledgementRangesChunk) Type() byte {
	return DataAcknowledgementRangesChunkType
}

func (chnk *DataAcknowledgementRangesChunk) Len() uint16 {
	l := 1 +
		chnk.FlowID.ByteLength() +
		chnk.BufferBlocksAvailable.ByteLength() +
		chnk.CumulativeAck.ByteLength()

	for _, r := range chnk.Ranges {
		l += r.HolesMinusOne.ByteLength() +
			r.ReceivedMinusOne.ByteLength()
	}

	return uint16(l)
}

func (chnk *DataAcknowledgementRangesChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	// Contents
	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return nil
	}

	if err = chnk.FlowID.WriteTo(buffer); err != nil {
		return nil
	}

	if err = chnk.BufferBlocksAvailable.WriteTo(buffer); err != nil {
		return nil
	}

	if err = chnk.CumulativeAck.WriteTo(buffer); err != nil {
		return nil
	}

	for i := range chnk.Ranges {
		chnk.Ranges[i].HolesMinusOne.WriteTo(buffer)
		chnk.Ranges[i].ReceivedMinusOne.WriteTo(buffer)
	}

	return nil
}

func (chnk *DataAcknowledgementRangesChunk) ReadFrom(buffer *bytes.Buffer) error {
	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	if err = chnk.FlowID.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.BufferBlocksAvailable.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.CumulativeAck.ReadFrom(buffer); err != nil {
		return err
	}

	rangesLen := int(length) -
		chnk.FlowID.ByteLength() -
		chnk.BufferBlocksAvailable.ByteLength() -
		chnk.CumulativeAck.ByteLength()

	l := list.New()

	i := rangesLen
	for i > 0 {
		holes := vlu.Vlu(0)
		recv := vlu.Vlu(0)

		if err = holes.ReadFrom(buffer); err != nil {
			break
		}

		i -= holes.ByteLength()

		if err = recv.ReadFrom(buffer); err != nil {
			break
		}

		i -= recv.ByteLength()

		if i < 0 {
			break
		}

		r := &DataAcknowledgementRange{
			HolesMinusOne:    holes,
			ReceivedMinusOne: recv,
		}

		l.PushBack(r)
	}

	if l.Len() > 0 {
		chnk.Ranges = make([]DataAcknowledgementRange, l.Len())

		i := 0
		for r := l.Front(); r != nil; r = r.Next() {
			chnk.Ranges[i] = *r.Value.(*DataAcknowledgementRange)
			i += 1
		}
	}

	return nil
}
