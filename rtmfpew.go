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
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/rtmfpew/rtmfpew/protocol"
)

const DefaultPort = 1935

func Connect(addr string) error {

	url, err := url.Parse(addr)
	if err != nil {
		return err
	}

	if url.Scheme != "rtmfp" {
		return errors.New("Protocol " + url.Scheme + " is not supported")
	}

	if len(url.Fragment) > 0 {
		log.Println("URL Fragment %s will be ignored", url.Fragment)
	}

	if len(url.Query().Encode()) > 0 {
		log.Println("URL Query %s will be ignored", url.Query().Encode())
	}

	if url.Host == ":" {
		url.Host = ""
	}

	if len(url.Host) > 0 {
		return rtmfpew.clientMode(url.Host)
	}

	return rtmfpew.serverlessMode()
}

func ListenSpecific(host string) error {
	return protocol.Run(host, protocol.ServerMode)
}

func Listen() error {
	return protocol.Run("0.0.0.0:"+strconv.Itoa(DefaultPort), protocol.ServerMode)
}
