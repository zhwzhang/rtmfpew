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

const BufferProbeChunkType = 0x18

// BufferProbeChunk is sent to request available receive buffer for a flow.
type BufferProbeChunk struct {
	FlowID io.Vlu
}

func (chnk *BufferProbeChunk) Type() byte {
	return BufferProbeChunkType
}

func (chnk *BufferProbeChunk) Len() uint16 {
	return 1 + uint16(chnk.FlowID.ByteLength())
}

func (chnk *BufferProbeChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	if err = chnk.FlowID.WriteTo(buffer); err != nil {
		return err
	}

	return nil
}

func (chnk *BufferProbeChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	if err = chnk.FlowID.ReadFrom(buffer); err != nil {
		return err
	}

	return nil
}
