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

func TestInitiatorInitialKeyingIO(t *testing.T) {
	Convey("Given a initiator initial keying chunk", t, func() {

		chnk := InitiatorInitialKeyingChunkSample()

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &InitiatorInitialKeyingChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, InitiatorInitialKeyingChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			for i := range chnk.CookieEcho {
				So(readChnk.CookieEcho[i], ShouldEqual, chnk.CookieEcho[i])
			}

			for i := range chnk.InitiatorCertificate {
				So(readChnk.InitiatorCertificate[i], ShouldEqual, chnk.InitiatorCertificate[i])
			}

			So(readChnk.InitiatorSessionID, ShouldEqual, chnk.InitiatorSessionID)

			for i := range chnk.SessionKeyInitiatorComponent {
				So(readChnk.SessionKeyInitiatorComponent[i], ShouldEqual, chnk.SessionKeyInitiatorComponent[i])
			}

			for i := range chnk.Signature {
				So(readChnk.Signature[i], ShouldEqual, chnk.Signature[i])
			}
		})
	})
}
