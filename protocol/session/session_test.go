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
	Convey("Given a packet with set of chunks", t, func() {
		sess := &Session{mtu: 20}
		pckt := &Packet{Chunks: list.New()}
		for i := 0; i < 3; i++ {
			pckt.Chunks.PushBack(chunks.InitiatorHelloChunkSample())
		}

		Convey("It should be fragmented properly", func() {
			chnks, err := pckt.doFragmentation(sess.mtu, &sess.pcktCounter)
			So(err, ShouldBeNil)
			So(chnks, ShouldNotBeNil)
			So(chnks.Len(), ShouldEqual, 2)
			fragmentedPckt := &Packet{Chunks: chnks}

			Convey("With properly sorted fragments", func() {
				i := uint32(0)
				for chnk := fragmentedPckt.Chunks.Front(); chnk != nil; chnk = chnk.Next() {
					packetID := uint32(chnk.Value.(*chunks.FragmentChunk).PacketID)
					frgNum := uint32(chnk.Value.(*chunks.FragmentChunk).FragmentNum)

					So(packetID, ShouldEqual, sess.pcktCounter)
					So(frgNum, ShouldEqual, i)

					i++
				}
			})

			Convey("And defragmented back", func() {
				buff := &fragmentsBuffer{}
				for chnk := fragmentedPckt.Chunks.Front(); chnk != nil; chnk = chnk.Next() {
					buff.Add(chnk.Value.(*chunks.FragmentChunk))
				}
				defragmentedPckt, err := sess.reassemblePacket(buff)

				So(err, ShouldBeNil)
				So(defragmentedPckt.TimeCritical, ShouldEqual, pckt.TimeCritical)
				So(defragmentedPckt.TimeCriticalReserve, ShouldEqual, pckt.TimeCriticalReserve)
				So(defragmentedPckt.Timestamp, ShouldEqual, pckt.Timestamp)
				So(defragmentedPckt.TimestampEcho, ShouldEqual, pckt.TimestampEcho)
				So(defragmentedPckt.TimestampPresent, ShouldEqual, pckt.TimestampPresent)
				So(defragmentedPckt.TimestampEchoPresent, ShouldEqual, pckt.TimestampEchoPresent)
				So(defragmentedPckt.Mode, ShouldEqual, pckt.Mode)
				So(defragmentedPckt.Chunks, ShouldResemble, pckt.Chunks)
			})
		})

		Convey("Packet should be written properly", func() {

		})

		Convey("And read back", func() {

		})
	})
}
