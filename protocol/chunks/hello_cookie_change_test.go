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
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestHelloCookieChangeIO(t *testing.T) {
	Convey("Given a hello cookie change chunk", t, func() {
		oldC := [...]byte{0x11, 0xBA, 0x2A, 0xEF, 0xA1}
		newC := [...]byte{0x12, 0x9A, 0x1A, 0xFD, 0x91}
		chnk := &HelloCookieChangeChunk{
			OldCookie: oldC[:],
			NewCookie: newC[:],
		}

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &HelloCookieChangeChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, HelloCookieChangeChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			for i := range chnk.OldCookie {
				So(readChnk.OldCookie[i], ShouldEqual, chnk.OldCookie[i])
			}

			for i := range chnk.NewCookie {
				So(readChnk.NewCookie[i], ShouldEqual, chnk.NewCookie[i])
			}
		})
	})
}