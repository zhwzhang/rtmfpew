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

const FragmentChunkType = 0x7f

// FragmentChunk is used to fragment packets in startup mode.
type FragmentChunk struct {
	MoreFragments bool
	PacketID      io.Vlu
	FragmentNum   io.Vlu

	Fragment []byte
}

// Type returns FragmentChunk type opcode.
func (chnk *FragmentChunk) Type() byte {
	return FragmentChunkType
}

func (chnk *FragmentChunk) Len() uint16 {
	return uint16(2 +
		chnk.PacketID.ByteLength() +
		chnk.FragmentNum.ByteLength() +
		len(chnk.Fragment))
}

func (chnk *FragmentChunk) WriteTo(buffer *bytes.Buffer) error {

	// Chunk header
	err := buffer.WriteByte(chnk.Type())
	if err != nil {
		return err
	}

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents
	if chnk.MoreFragments {
		err = buffer.WriteByte(128) // 7 bits are reserved
	} else {
		err = buffer.WriteByte(0)
	}

	if err != nil {
		return err
	}

	if err = chnk.PacketID.WriteTo(buffer); err != nil {
		return err
	}

	if err = chnk.FragmentNum.WriteTo(buffer); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.Fragment); err != nil {
		return err
	}

	return nil
}

func (chnk *FragmentChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	moreFragments, err := buffer.ReadByte()
	if err != nil {
		return err
	}

	chnk.MoreFragments = (moreFragments == 128)

	if err = chnk.PacketID.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.FragmentNum.ReadFrom(buffer); err != nil {
		return err
	}

	fragmentLength := (int(length) - chnk.PacketID.ByteLength() - chnk.FragmentNum.ByteLength())
	chnk.Fragment = make([]byte, fragmentLength)

	if _, err = buffer.Read(chnk.Fragment); err != nil {
		return err
	}

	return nil
}
