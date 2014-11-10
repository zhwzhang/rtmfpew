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

const FlowExceptionReportChunkType = 0x5e

// FlowExceptionReportChunk is sent to close the flow.
type FlowExceptionReportChunk struct {
	FlowID    io.Vlu
	Exception io.Vlu
}

// Type returns FlowExceptionReportChunk type opcode
func (chnk *FlowExceptionReportChunk) Type() byte {
	return FlowExceptionReportChunkType
}

func (chnk *FlowExceptionReportChunk) Len() uint16 {
	return uint16(1 + chnk.FlowID.ByteLength() +
		chnk.Exception.ByteLength())
}

func (chnk *FlowExceptionReportChunk) WriteTo(buffer *bytes.Buffer) error {

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

	if err = chnk.Exception.WriteTo(buffer); err != nil {
		return err
	}

	return nil
}

func (chnk *FlowExceptionReportChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	if err = chnk.FlowID.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.Exception.ReadFrom(buffer); err != nil {
		return err
	}

	return nil
}
