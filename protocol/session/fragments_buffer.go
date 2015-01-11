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

import "github.com/rtmfpew/rtmfpew/protocol/chunks"

type fragmentsBuffer struct {
	Size uint16
	buff []*chunks.FragmentChunk
}

func CreateNewFragmentsBuffer(chnk *chunks.FragmentChunk) *fragmentsBuffer {
	f := &fragmentsBuffer{}
	f.Add(chnk)
	return f
}

func (f *fragmentsBuffer) Len() int {
	return len(f.buff)
}

func (f *fragmentsBuffer) Add(chunk *chunks.FragmentChunk) {
	num := int(chunk.FragmentNum)
	switch {
	case num+1 > len(f.buff):
		old := f.buff
		f.buff = make([]*chunks.FragmentChunk, num+1)
		if old != nil {
			copy(f.buff, old)
		}
	case f.buff[num] != nil:
		return
	}
	f.buff[num] = chunk
	f.Size += chunk.Len()
}

func (f *fragmentsBuffer) IsComplete() bool {
	last := len(f.buff) - 1
	if f.buff[last] == nil || f.buff[last].MoreFragments {
		return false
	}
	for i := last - 1; i >= 0; i-- {
		if f.buff[i] == nil {
			return false
		}
	}
	return true
}
