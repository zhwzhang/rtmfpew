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

func TestFragmentIO(t *testing.T) {
	Convey("Given a fragment chunk", t, func() {

		chnk := FragmentChunkSample()

		buff := bytes.NewBuffer(make([]byte, 0))

		err := chnk.WriteTo(buff)
		Convey("It can be written into a buffer", func() {
			So(err, ShouldBeNil)
		})

		Convey("It can be read back from the buffer", func() {
			readChnk := &FragmentChunk{}
			typ, err := buff.ReadByte()

			So(err, ShouldBeNil)
			So(typ, ShouldEqual, FragmentChunkType)

			err = readChnk.ReadFrom(buff)
			So(err, ShouldBeNil)

			So(readChnk.MoreFragments, ShouldEqual, chnk.MoreFragments)
			So(uint32(readChnk.PacketID), ShouldEqual, uint32(chnk.PacketID))
			So(uint32(readChnk.FragmentNum), ShouldEqual, uint32(chnk.FragmentNum))

			for i := range chnk.Fragment {
				So(readChnk.Fragment[i], ShouldEqual, chnk.Fragment[i])
			}
		})
	})
}
