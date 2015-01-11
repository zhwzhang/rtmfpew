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
	"testing"

	"github.com/rtmfpew/rtmfpew/protocol/chunks"

	"github.com/rtmfpew/amfy/vlu"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFragmentsBuffer(t *testing.T) {
	Convey("Given an uncomplete set of fragment chunks", t, func() {
		b := &fragmentsBuffer{}
		ch1 := chunks.FragmentChunkSample()
		ch1.FragmentNum = 5
		b.Add(ch1)
		ch2 := chunks.FragmentChunkSample()
		ch2.FragmentNum = 0
		b.Add(ch2)
		Convey("They should be added properly", func() {
			So(b.buff[5], ShouldEqual, ch1)
			So(b.buff[0], ShouldEqual, ch2)
			So(b.Size, ShouldEqual, ch1.Len()+ch2.Len())
			So(b.Len(), ShouldEqual, 6)
		})
		Convey("Buffer should not be completed", func() {
			So(b.IsComplete(), ShouldBeFalse)
		})
	})

	Convey("And with complete set of fragment chunks", t, func() {
		b := &fragmentsBuffer{}
		var ch *chunks.FragmentChunk
		for i := vlu.Vlu(0); i < 10; i++ {
			ch = chunks.FragmentChunkSample()
			ch.FragmentNum = i
			b.Add(ch)
		}
		ch.MoreFragments = false
		Convey("Buffer should be completed", func() {
			So(b.IsComplete(), ShouldBeTrue)
		})
	})
}
