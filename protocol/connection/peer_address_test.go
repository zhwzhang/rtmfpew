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

package connection

import (
	"bytes"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPeerAddressIO(t *testing.T) {
	Convey("Given two buffers with IPv4 and IPv6 peer addresses", t, func() {

		buff4 := bytes.NewBuffer(make([]byte, 0))
		buff16 := bytes.NewBuffer(make([]byte, 0))
		testAddr4, _ := net.ResolveUDPAddr("udp", "53.13.1.45:1935")
		testAddr16, _ := net.ResolveUDPAddr("udp", "53.13.1.45:1935")

		addr4 := &PeerAddress{
			IP:     []byte(testAddr4.IP.To4()),
			Port:   uint16(testAddr4.Port),
			Origin: LocalOrigin,
		}

		addr16 := &PeerAddress{
			IP:     []byte(testAddr16.IP.To16()),
			Port:   uint16(testAddr16.Port),
			Origin: RemoteOrigin,
		}

		buff4.Reset()
		buff16.Reset()

		err := addr4.WriteTo(buff4)
		So(err, ShouldBeNil)
		err = addr16.WriteTo(buff16)
		So(err, ShouldBeNil)

		Convey("Should read IPv4 address", func() {
			newAddr4 := &PeerAddress{}
			newAddr4.ReadFrom(buff4)

			for i := range newAddr4.IP {
				So(newAddr4.IP[i], ShouldEqual, addr4.IP[i])
			}

			So(newAddr4.Port, ShouldEqual, addr4.Port)
			So(newAddr4.Origin, ShouldEqual, addr4.Origin)
		})

		Convey("Should read IPv6 address", func() {
			newAddr16 := &PeerAddress{}
			newAddr16.ReadFrom(buff16)

			for i := range newAddr16.IP {
				So(newAddr16.IP[i], ShouldEqual, addr16.IP[i])
			}

			So(newAddr16.Port, ShouldEqual, addr16.Port)
			So(newAddr16.Origin, ShouldEqual, addr16.Origin)
		})
	})
}
