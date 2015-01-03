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
	"container/list"
	"testing"

	"github.com/rtmfpew/rtmfpew/protocol/chunks"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSessionChunksFragmentation(t *testing.T) {
	Convey("Given a set of chunks", t, func() {
		chnks := list.New()
		for i := 0; i < 3; i++ {
			chnks.PushBack(chunks.InitiatorHelloChunkSample()) // len is 12
		}
		fitMTU := (&Session{mtu: 100}).fragmentChunks(chnks)
		sess := &Session{mtu: 20}
		largerThanMTU := sess.fragmentChunks(chnks)

		Convey("They should be fragmented properly", func() {
			So(largerThanMTU.Len(), ShouldEqual, 2)
			So(fitMTU, ShouldResemble, chnks)
		})

		Convey("Fragments should be sorted", func() {
			i := uint32(0)
			for chnk := largerThanMTU.Front(); chnk != nil; chnk = chnk.Next() {
				packetID := uint32(chnk.Value.(*chunks.FragmentChunk).PacketID)
				frgNum := uint32(chnk.Value.(*chunks.FragmentChunk).FragmentNum)

				So(packetID, ShouldEqual, sess.pcktCounter)
				So(frgNum, ShouldEqual, i)

				i++
			}
		})

		Convey("And defragmented back", func() {

		})

		Convey("Packet should be written properly", func() {

		})

		Convey("And read back", func() {

		})
	})
}
