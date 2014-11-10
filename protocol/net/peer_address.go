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

package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/rtmfpew/rtmfpew/protocol/io"
	"net"
)

// Address origins
const (
	UnknownOrigin = 0
	LocalOrigin
	RemoteOrigin
	ProxyOrigin
)

// PeerAddress is a simple
type PeerAddress struct {
	IP     []byte
	Port   uint16
	Origin byte
}

// Length returns PeerAddress length
func (addr *PeerAddress) Length() int {
	return 1 + // Origin
		len(addr.IP) + // IP
		2 // port uint16
}

func PeerAddressFrom(udpAddr *net.UDPAddr) *PeerAddress {
	addr := &PeerAddress{
		IP: []byte(udpAddr.IP),
		Origin: LocalOrigin,
		Port: uint16(udpAddr.Port),
	}
	
	return addr
}

func (addr *PeerAddress) ReadFrom(buffer *bytes.Buffer) (err error) {

	flags, err := buffer.ReadByte()
	if err != nil {
		return err
	}

	addr.Origin = flags & 3 // first two bits, 5 reserved

	if io.BitIsSet(&flags, 7) { // is IPv6
		addr.IP = make([]byte, net.IPv6len)
	} else {
		addr.IP = make([]byte, net.IPv4len)
	}

	num, err := buffer.Read(addr.IP)
	if err != nil {
		return err
	}

	if num != net.IPv4len && num != net.IPv6len {
		return errors.New("Can't read IP addr")
	}

	if err = binary.Read(buffer, binary.BigEndian, &addr.Port); err != nil {
		return err
	}

	return nil
}

func (addr *PeerAddress) WriteTo(buffer *bytes.Buffer) error {

	flags := addr.Origin
	if len(addr.IP) == net.IPv6len {
		flags |= (1 << 7)
	}

	err := buffer.WriteByte(flags)
	if err != nil {
		return err
	}

	num, err := buffer.Write(addr.IP)
	if err != nil {
		return err
	}

	if num < len(addr.IP) {
		return errors.New("Can't write IP addr")
	}

	if err = binary.Write(buffer, binary.BigEndian, addr.Port); err != nil {
		return err
	}

	return nil
}
