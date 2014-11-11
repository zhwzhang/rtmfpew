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

package config

import (
// "github.com/rtmfpew/rtmfpew/protocol/vlu"
)

type configValues struct {
	Mtu                 uint
	MaxFragmentationGap uint
	MaxFragments        uint
}

var values = &configValues{
	Mtu:                 768,
	MaxFragmentationGap: 3,
	MaxFragments:        4,
}

// Load loads config values from file
func Load() {

}

// Mtu returens packet mtu value
func Mtu() uint {
	return values.Mtu
}

// MaxFragmentationGap returns max gap for fragmented packets
func MaxFragmentationGap() uint {
	return values.MaxFragmentationGap
}

// MaxFragments returns max number of packet fragments
func MaxFragments() uint {
	return values.MaxFragments
}
