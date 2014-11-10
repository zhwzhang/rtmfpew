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
)

const NextUserDataChunkType = 0x11

// NextUserDataChunk equivalent to UserData
type NextUserDataChunk UserDataChunk

// Type returns NextUserDataChunk type opcode.
func (chnk *NextUserDataChunk) Type() byte {
	return NextUserDataChunkType
}

func (chnk *NextUserDataChunk) Len() uint16 {
	dataChnk := UserDataChunk(*chnk)
	return dataChnk.Len()
}

func (chnk *NextUserDataChunk) WriteTo(buffer *bytes.Buffer) error {
	dataChnk := UserDataChunk(*chnk)
	err := dataChnk.WriteNextTo(buffer, chnk.Type())
	if err != nil {
		return err
	}

	return nil
}

func (chnk *NextUserDataChunk) ReadFrom(buffer *bytes.Buffer) error {
	dataChnk := UserDataChunk(*chnk)
	err := dataChnk.ReadFrom(buffer)
	if err != nil {
		return err
	}

	return nil
}
