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
	"errors"
	"github.com/rtmfpew/amfy/vlu"
)

const UserDataChunkType = 0x10

// UserDataChunk fragment control modes
const (
	WholeFragmentControl = 0
	BeginFragmentControl
	EndFragmentControl
	MiddleFragmentControl
)

// UserDataChunk is a basic unit of transmission for the user messages of a flow
// Valid in established session and Initiator / Responder packet modes
type UserDataChunk struct {
	OptionsPresent  bool
	FragmentControl byte
	Abandon         bool
	Final           bool
	FlowID          vlu.Vlu
	SequenceNumber  vlu.Vlu
	FsnOffset       vlu.Vlu

	Options  []UserDataOption
	UserData []byte
}

// Type returns UserDataChunk type opcode
func (chnk *UserDataChunk) Type() byte {
	return UserDataChunkType
}

func (chnk *UserDataChunk) Len() uint16 {
	l := 1 +
		1 + // flags
		chnk.FlowID.ByteLength() +
		chnk.SequenceNumber.ByteLength() +
		chnk.FsnOffset.ByteLength() +
		len(chnk.UserData)
	if chnk.OptionsPresent {
		for _, opt := range chnk.Options {
			l += opt.Length()
		}

		l += 1 // for opt list marker
	}

	return uint16(l)
}

func (chnk *UserDataChunk) WriteTo(buffer *bytes.Buffer) error {
	return chnk.WriteNextTo(buffer, chnk.Type())
}

func (chnk *UserDataChunk) WriteNextTo(buffer *bytes.Buffer, typ byte) error {

	// Chunk header
	err := buffer.WriteByte(typ)
	if err != nil {
		return err
	}

	chnk.OptionsPresent = len(chnk.Options) > 0

	if err = binary.Write(buffer, binary.BigEndian, chnk.Len()-1); err != nil {
		return err
	}

	// Contents

	// Pack Flags
	flags := byte(0)
	if chnk.OptionsPresent {
		flags++
	}

	flags <<= 3 // one bit reserved
	flags += chnk.FragmentControl
	flags <<= 3 // two bits reserved
	if chnk.Abandon {
		flags++
	}
	flags <<= 1
	if chnk.Final {
		flags++
	}

	if err = buffer.WriteByte(flags); err != nil {
		return err
	}

	if err = chnk.FlowID.WriteTo(buffer); err != nil {
		return err
	}

	if err = chnk.SequenceNumber.WriteTo(buffer); err != nil {
		return err
	}

	if err = chnk.FsnOffset.WriteTo(buffer); err != nil {
		return err
	}

	if chnk.OptionsPresent {
		for _, opt := range chnk.Options {
			if err = opt.WriteTo(buffer); err != nil {
				return err
			}
		}
	}

	marker := vlu.Vlu(0)
	if err = marker.WriteTo(buffer); err != nil {
		return err
	}

	if _, err = buffer.Write(chnk.UserData); err != nil {
		return err
	}

	return nil
}

func (chnk *UserDataChunk) ReadFrom(buffer *bytes.Buffer) error {

	// Chunk header
	length := uint16(0)
	err := binary.Read(buffer, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	// Contents
	flags, err := buffer.ReadByte()
	if err != nil {
		return err
	}

	// Unpack Flags
	chnk.Final = (flags & 1) != 0
	flags >>= 1
	chnk.Abandon = (flags & 1) != 0
	flags >>= 3
	chnk.FragmentControl = (flags & 3)
	flags >>= 3
	chnk.OptionsPresent = (flags & 1) != 0

	if err = chnk.FlowID.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.SequenceNumber.ReadFrom(buffer); err != nil {
		return err
	}

	if err = chnk.FsnOffset.ReadFrom(buffer); err != nil {
		return err
	}

	dataLength := int(length) - 1 -
		chnk.FlowID.ByteLength() -
		chnk.SequenceNumber.ByteLength() -
		chnk.FsnOffset.ByteLength()

	if chnk.OptionsPresent {
		optList := list.New()

		optLen := vlu.Vlu(1)
		for dataLength > 0 {
			opt := UserDataOption{}
			optLen, _ = opt.ReadFrom(buffer)

			if optLen == 0 {
				break
			}

			dataLength -= opt.Length()
			optList.PushBack(opt)
		}

		if dataLength <= 0 {
			return errors.New("Corrupted data packet")
		}

		chnk.Options = make([]UserDataOption, optList.Len())

		i := 0
		for opt := optList.Front(); opt != nil; opt = opt.Next() {
			chnk.Options[i] = opt.Value.(UserDataOption)
			i += 1
		}
	}

	chnk.UserData = make([]byte, dataLength)
	if _, err := buffer.Read(chnk.UserData); err != nil {
		return err
	}

	return nil
}
