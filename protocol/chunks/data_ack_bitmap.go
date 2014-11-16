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

const DataAcknowledgementBitmapChunkType = 0x50

// DataAcknowledgementBitmapChunk is sent to indicate UserData fragment sequense numbers beeing recieved.
type DataAcknowledgementBitmapChunk struct {
	FlowID                vlu.Vlu
	BufferBlocksAvailable vlu.Vlu
	CumulativeAck         vlu.Vlu

	Acknowledgement []byte
}

// Type returns DataAcknowledgementBitmapChunk type opcode.
func (chnk *DataAcknowledgementBitmapChunk) Type() byte {
	return DataAcknowledgementBitmapChunkType
}

func (chnk *DataAcknowledgementBitmapChunk) Len() uint16 {
	return uint16(1 +
		chnk.FlowID.ByteLength() +
		chnk.BufferBlocksAvailable.ByteLength() +
		chnk.CumulativeAck.ByteLength() +
		int(len(chnk.Acknowledgement)))
}

func (chnk *DataAcknowledgementBitmapChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents

	if err = chnk.FlowID.WriteTo(buffer); err != nil {
		return err
	}

	if err = chnk.BufferBlocksAvailable.WriteTo(buffer); err != nil {
		return err
	}

	if err = chnk.CumulativeAck.WriteTo(buffer); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Acknowledgement); err != nil {
		return err
	}

	return nil
}

func (chnk *DataAcknowledgementBitmapChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	if err = chnk.FlowID.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.BufferBlocksAvailable.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.CumulativeAck.ReadFrom(buffer); err != nil {
		return err
	}

	ackLength := int(length) -
		chnk.FlowID.ByteLength() -
		chnk.BufferBlocksAvailable.ByteLength() -
		chnk.CumulativeAck.ByteLength()

	chnk.Acknowledgement = make([]byte, ackLength)
	num, err := buffer.Read(chnk.Acknowledgement)

	if err != nil {
		return err
	}

	if num < ackLength {
		return errors.New("Can't read DataAckBitmap chunk acknowledgement")
	}

	return nil
}
