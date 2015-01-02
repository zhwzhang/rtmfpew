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

package ip

import (
	// "bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChecksum(t *testing.T) {

	data := []byte{
		0x45, 0x00, 0x00, 0x47,
		0x73, 0x88, 0x40, 0x00,
		0x40, 0x06, 0xA2, 0xC4,
		0x83, 0x9F, 0x0E, 0x85,
		0x83, 0x9F, 0x0E, 0xA1,
	}

	Convey("Given a test data chunk", t, func() {
		Convey("Checksum should be calculated properly", func() {
			So(Checksum(data), ShouldEqual, 0x0000)
		})
	})
}
