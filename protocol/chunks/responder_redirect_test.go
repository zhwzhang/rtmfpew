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
	. "github.com/rtmfpew/rtmfpew/protocol/net"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResponderRedirectIO(t *testing.T) {
	Convey("Given a responder redirect chunk", t, func() {
		t := [...]byte{0x2A, 0xC3, 0xB1, 0x5C}
		a := [...]PeerAddress{
			PeerAddress{
				IP:     t[:],
				Port:   2913,
				Origin: RemoteOrigin,
			}, PeerAddress{
				IP:     t[:],
				Port:   2911,
				Origin: ProxyOrigin,
			}, PeerAddress{
				IP:     t[:],
				Port:   2912,
				Origin: LocalOrigin,
			},
		}

		chnk := &ResponderRedirectChunk{
			RedirectDestination: a[:],
			TagEcho:             t[:],
		}

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &ResponderRedirectChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, ResponderRedirectChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			for i := range chnk.RedirectDestination {
				for j := range chnk.RedirectDestination[i].IP {
					So(readChnk.RedirectDestination[i].IP[j], ShouldEqual, readChnk.RedirectDestination[i].IP[j])
				}

				So(readChnk.RedirectDestination[i].Origin, ShouldEqual, readChnk.RedirectDestination[i].Origin)
				So(readChnk.RedirectDestination[i].Port, ShouldEqual, readChnk.RedirectDestination[i].Port)
			}

			for i := range chnk.TagEcho {
				So(readChnk.TagEcho[i], ShouldEqual, readChnk.TagEcho[i])
			}
		})
	})
}
