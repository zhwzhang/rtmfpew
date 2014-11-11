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

func TestUserDataIO(t *testing.T) {
	Convey("Given a user data chunk", t, func() {

		chnk := UserDataChunkSample()

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {

			readChnk := &UserDataChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, UserDataChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			So(readChnk.Abandon, ShouldEqual, chnk.Abandon)
			So(readChnk.Final, ShouldEqual, chnk.Final)
			So(readChnk.FragmentControl, ShouldEqual, chnk.FragmentControl)
			So(readChnk.OptionsPresent, ShouldEqual, chnk.OptionsPresent)

			So(uint32(readChnk.FlowID), ShouldEqual, uint32(chnk.FlowID))
			So(uint32(readChnk.FsnOffset), ShouldEqual, uint32(chnk.FsnOffset))
			So(uint32(readChnk.SequenceNumber), ShouldEqual, uint32(chnk.SequenceNumber))

			for i := range chnk.Options {
				So(readChnk.Options[i].OptionType, ShouldEqual, chnk.Options[i].OptionType)
				for j := range chnk.Options[i].Value {
					So(readChnk.Options[i].Value[j], ShouldEqual, chnk.Options[i].Value[j])
				}
			}

			for i := range chnk.UserData {
				So(readChnk.UserData[i], ShouldEqual, chnk.UserData[i])
			}
		})
	})
}
