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

package rtmfp

import (
	"net"

	"github.com/rtmfpew/rtmfpew/protocol"
)

func serverlessMode() (*protocol.Context, error) {

	addr := net.ResolveUDPAddr("udp", "0.0.0.0")
	if addr.Port == 0 {
		addr.Port = DefaultPort
	}

	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	ctx := &protocol.Context{}

	return ctx, nil
}

func clientMode(host string) (*protocol.Context, error) {

	addr := net.ResolveUDPAddr("udp", host)
	if addr.Port == 0 {
		addr.Port = DefaultPort
	}

	conn, err := net.DialUDP("udp", nil, addr)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	ctx := &protocol.Context{}

	return ctx, nil
}

func serverMode(host string) (*protocol.Context, error) {

	addr := net.ResolveUDPAddr("udp", host)
	if addr.Port == 0 {
		addr.Port = DefaultPort
	}

	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	ctx := &protocol.Context{}

	return ctx, nil
}
