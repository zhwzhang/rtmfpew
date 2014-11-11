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
	"github.com/rtmfpew/rtmfpew/protocol/vlu"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDataAckRangesIO(t *testing.T) {
	Convey("Given a data ack bitmap chunk", t, func() {
		ranges := [...]DataAcknowledgementRange{
			DataAcknowledgementRange{
				HolesMinusOne:    vlu.Vlu(141),
				ReceivedMinusOne: vlu.Vlu(1492351),
			}, DataAcknowledgementRange{
				HolesMinusOne:    vlu.Vlu(1112),
				ReceivedMinusOne: vlu.Vlu(119),
			}, DataAcknowledgementRange{
				HolesMinusOne:    vlu.Vlu(10131),
				ReceivedMinusOne: vlu.Vlu(109),
			}, DataAcknowledgementRange{
				HolesMinusOne:    vlu.Vlu(151),
				ReceivedMinusOne: vlu.Vlu(1191),
			},
		}

		chnk := &DataAcknowledgementRangesChunk{
			FlowID:                vlu.Vlu(1831),
			CumulativeAck:         vlu.Vlu(119),
			BufferBlocksAvailable: vlu.Vlu(284),
			Ranges:                ranges[:],
		}

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &DataAcknowledgementRangesChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, DataAcknowledgementRangesChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			So(uint32(readChnk.FlowID), ShouldEqual, uint32(chnk.FlowID))
			So(uint32(readChnk.CumulativeAck), ShouldEqual, uint32(chnk.CumulativeAck))
			So(uint32(readChnk.BufferBlocksAvailable), ShouldEqual, uint32(chnk.BufferBlocksAvailable))

			for i := range chnk.Ranges {
				So(uint32(readChnk.Ranges[i].HolesMinusOne), ShouldEqual, uint32(chnk.Ranges[i].HolesMinusOne))
				So(uint32(readChnk.Ranges[i].ReceivedMinusOne), ShouldEqual, uint32(chnk.Ranges[i].ReceivedMinusOne))
			}
		})
	})
}
