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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestForwardedHelloIO(t *testing.T) {
	Convey("Given a forworded hello chunk", t, func() {

		chnk := ForwardedHelloChunkSample()

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &ForwardedHelloChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, ForwardedHelloChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			for i := range chnk.Epd {
				So(readChnk.Epd[i], ShouldEqual, chnk.Epd[i])
			}

			for i := range chnk.Tag {
				So(readChnk.Tag[i], ShouldEqual, chnk.Tag[i])
			}

			for i := range chnk.ReplyAddress.IP {
				So(readChnk.ReplyAddress.IP[i], ShouldEqual, chnk.ReplyAddress.IP[i])
			}

			So(readChnk.ReplyAddress.Origin, ShouldEqual, chnk.ReplyAddress.Origin)
			So(readChnk.ReplyAddress.Port, ShouldEqual, chnk.ReplyAddress.Port)
		})
	})
}
