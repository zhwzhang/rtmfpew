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

package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// Profile rtmfp encryption profile interface
type Profile interface {
	Init(key []byte) error
	InitDefault() error
	Encrypt(data bytes.Buffer) error
	Decrypt(data bytes.Buffer) error
}

// DefaultProfile for rmtmfp encryption
type DefaultProfile struct {
	key         []byte
	blockCipher cipher.Block
}

// DefaultKey for rtmfp encryption
var DefaultKey = [...]byte{ // Adobe Systems 02
	0x41, 0x64, 0x6F, 0x62,
	0x65, 0x20, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6D,
	0x73, 0x20, 0x30, 0x32,
}

// Init crypto profile with specific encryption key
func (profile *DefaultProfile) Init(key []byte) error {
	var err error

	if len(key) > 0 {
		profile.blockCipher, err = aes.NewCipher(key)

		profile.key = key

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Encryption key required")
}

// InitDefault init crypto profile with default encryption key
func (profile *DefaultProfile) InitDefault() error {
	return profile.Init(DefaultKey[:16])
}

// Encrypt encrupts buffer with specified block cipher
func (profile *DefaultProfile) Encrypt(data bytes.Buffer) error {
	if profile.blockCipher != nil {
		encryptedBuf := make([]byte, len(data.Bytes()))

		profile.blockCipher.Encrypt(encryptedBuf, data.Bytes())
		data.Reset()
		data.Write(encryptedBuf)

		return nil
	}

	return errors.New("Init crypto profile first")
}

// Decrypt decrypts buffer with specified block cipher
func (profile *DefaultProfile) Decrypt(data bytes.Buffer) error {
	if profile.blockCipher != nil {
		decryptedBuf := make([]byte, len(data.Bytes()))

		profile.blockCipher.Decrypt(decryptedBuf, data.Bytes())
		data.Reset()
		data.Write(decryptedBuf)

		return nil
	}

	return errors.New("Init crypto profile first")
}
