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
	"errors"
	"github.com/rtmfpew/rtmfpew/protocol/vlu"
)

// Default UserDataOption types
const (
	PerFlowMetadataOptionType       = 0x00
	ReturnFlowAssociationOptionType = 0x0a
)

// UserDataOption is a container of option values.
type UserDataOption struct {
	OptionType vlu.Vlu
	Value      []byte
}

func (opt *UserDataOption) Length() int {
	l := vlu.Vlu(opt.OptionType.ByteLength() + len(opt.Value))
	return l.ByteLength() + int(l)
}

func (opt *UserDataOption) WriteTo(buffer *bytes.Buffer) error {
	length := vlu.Vlu(opt.OptionType.ByteLength() + len(opt.Value))
	err := length.WriteTo(buffer)
	if err != nil {
		return err
	}

	if err = opt.OptionType.WriteTo(buffer); err != nil {
		return err
	}

	if _, err = buffer.Write(opt.Value); err != nil {
		return err
	}

	return nil
}

func (opt *UserDataOption) ReadFrom(buffer *bytes.Buffer) (vlu.Vlu, error) {

	length := vlu.Vlu(0)
	err := length.ReadFrom(buffer)
	if err != nil {
		return 0, err
	}

	if int(length) == 0 { // It's a Marker
		return 0, nil	
	}

	if err = opt.OptionType.ReadFrom(buffer); err != nil {
		return 0, err
	}

	valueLength := int(length) - opt.OptionType.ByteLength()

	opt.Value = make([]byte, valueLength)
	num, err := buffer.Read(opt.Value)

	if err != nil {
		return 0, err
	}

	if num < valueLength {
		return 0, errors.New("Can't read user data option value")
	}

	return length, nil
}
