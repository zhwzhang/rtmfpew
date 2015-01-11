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

func TestUserDataOptionIO(t *testing.T) {
	Convey("Given a user data option chunk", t, func() {

		chnk := UserDataOptionSample()

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &UserDataOption{}

			_, err := readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			So(uint32(readChnk.OptionType), ShouldEqual, uint32(chnk.OptionType))
			for i := range chnk.Value {
				So(readChnk.Value[i], ShouldEqual, chnk.Value[i])
			}
		})
	})
}
