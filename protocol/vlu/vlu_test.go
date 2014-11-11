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

package vlu

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVluIO(t *testing.T) {
	Convey("Given a set of Vlu's with different values and byte buffers", t, func() {
		oneByteVlu := Vlu(117)
		twoByteVlu := Vlu(14180)
		threeByteVlu := Vlu(1986010)
		fourByteVlu := Vlu(218231350)

		Convey("Bit operators should perform properly", func() {
			b := byte(0)
			SetBit(&b, 7)
			SetBit(&b, 6)
			So(b, ShouldEqual, 192)

			So(BitIsSet(&b, 7), ShouldBeTrue)
			So(BitIsSet(&b, 6), ShouldBeTrue)
			So(BitIsSet(&b, 5), ShouldBeFalse)
		})

		Convey("Vlu's should return proper byte lengths", func() {
			So(oneByteVlu.ByteLength(), ShouldEqual, 1)
			So(twoByteVlu.ByteLength(), ShouldEqual, 2)
			So(threeByteVlu.ByteLength(), ShouldEqual, 3)
			So(fourByteVlu.ByteLength(), ShouldEqual, 4)
		})

		Convey("Should be able to get 7-bit values", func() {
			So(get7BitValue(uint32(oneByteVlu), 0), ShouldEqual, 0x75)
			So(get7BitValue(uint32(twoByteVlu), 0), ShouldEqual, 0x64)
			So(get7BitValue(uint32(twoByteVlu), 1), ShouldEqual, 0x6E)
		})

		oneByteSlice := make([]byte, 1)
		twoByteSlice := make([]byte, 2)
		threeByteSlice := make([]byte, 3)
		fourByteSlice := make([]byte, 4)

		oneByteBuff := bytes.NewBuffer(oneByteSlice)
		twoByteBuff := bytes.NewBuffer(twoByteSlice)
		threeByteBuff := bytes.NewBuffer(threeByteSlice)
		fourByteBuff := bytes.NewBuffer(fourByteSlice)

		oneByteBuff.Reset()
		twoByteBuff.Reset()
		threeByteBuff.Reset()
		fourByteBuff.Reset()

		oneByteVlu.WriteTo(oneByteBuff)
		twoByteVlu.WriteTo(twoByteBuff)
		threeByteVlu.WriteTo(threeByteBuff)
		fourByteVlu.WriteTo(fourByteBuff)

		Convey("Vlu's should be written to buffer with proper lengths", func() {
			So(oneByteBuff.Len(), ShouldEqual, 1)
			So(twoByteBuff.Len(), ShouldEqual, 2)
			So(threeByteBuff.Len(), ShouldEqual, 3)
			So(fourByteBuff.Len(), ShouldEqual, 4)
		})

		actualOneByteSlice := [...]byte{0x75}
		actualTwoByteSlice := [...]byte{0xEE, 0x64}
		Convey("Vlu's should be properly written", func() {
			So(oneByteSlice[0], ShouldEqual, actualOneByteSlice[0])
			So(twoByteSlice[0], ShouldEqual, actualTwoByteSlice[0])
			So(twoByteSlice[1], ShouldEqual, actualTwoByteSlice[1])
		})

		readOneByteVlu := Vlu(0)
		readTwoByteVlu := Vlu(0)
		readThreeByteVlu := Vlu(0)
		readFourByteVlu := Vlu(0)

		Convey("Should be able to be read Vlu's back from buffers", func() {

			readOneByteVlu.ReadFrom(oneByteBuff)
			readTwoByteVlu.ReadFrom(twoByteBuff)
			readThreeByteVlu.ReadFrom(threeByteBuff)
			readFourByteVlu.ReadFrom(fourByteBuff)

			So(uint32(readOneByteVlu), ShouldEqual, uint32(oneByteVlu))
			So(uint32(readTwoByteVlu), ShouldEqual, uint32(twoByteVlu))
			So(uint32(readThreeByteVlu), ShouldEqual, uint32(threeByteVlu))
			So(uint32(readFourByteVlu), ShouldEqual, uint32(fourByteVlu))
		})
	})
}
