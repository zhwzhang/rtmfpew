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

package session

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"io"
	"sync/atomic"

	"github.com/rtmfpew/rtmfpew/config"
	"github.com/rtmfpew/rtmfpew/protocol/chunks"
	"github.com/rtmfpew/rtmfpew/protocol/connection"
	"github.com/rtmfpew/rtmfpew/protocol/crypto"
	"github.com/rtmfpew/rtmfpew/protocol/ip"

	"github.com/rtmfpew/amfy/vlu"
)

var (
	packetMtu           = config.Mtu()
	maxFragmentationGap = config.MaxFragmentationGap()
	maxFragments        = config.MaxFragments()
	maxFragmentsSize    = config.MaxFragmentsSize()
)

// SessionType stores current session state and validates it's changes
type SessionType interface {
	IsValidChunkType() uint8
	GotChunkType(uint8)
	NextType() *SessionType
}

// Session stores current connection state
type Session struct {
	ID uint32

	InitiatorAddr *connection.PeerAddress
	ResponderAddr *connection.PeerAddress

	profile crypto.Profile

	pcktCounter uint32
	mtu         uint16

	HasChecksums bool
	Established  bool

	fragmented map[vlu.Vlu]*fragmentsBuffer

	fragments     map[vlu.Vlu]*list.List
	fragmentSizes map[vlu.Vlu]uint16

	Type SessionType
}

// NewWith creates new session with custom profile
func NewWith(profile crypto.Profile, t SessionType) *Session {
	session := &Session{
		profile:      profile,
		pcktCounter:  0,
		mtu:          uint16(packetMtu),
		HasChecksums: false,
		Established:  false,
		Type:         t,
	}

	return session
}

// New creates new session with default profile
func New(t SessionType) *Session {
	profile := &crypto.DefaultProfile{}
	profile.InitDefault()

	return NewWith(profile, t)
}

// SetEncryptionKey for data cipher
func (session *Session) SetEncryptionKey(key []byte) error {
	return session.profile.Init(key)
}

func (session *Session) encryptBuffer(buff *bytes.Buffer) error {
	return session.profile.EncryptAt(buff, 4) // ID size
}

func (session *Session) decryptBuffer(buff *bytes.Buffer) error {
	return session.profile.DecryptAt(buff, 4) // ID size
}

// readID reads session ID
func (session *Session) readID(buff io.Reader) error {
	err := binary.Read(buff, binary.BigEndian, &session.ID)
	if err != nil {
		return err
	}

	b, err := bufio.NewReader(buff).Peek(2) // Unscramble session ID
	session.ID = session.ID ^ uint32(b[0]) ^ uint32(b[1])

	return nil
}

// writeID writes session ID
func (session *Session) writeID(buff *bytes.Buffer) error {
	buff = bytes.NewBuffer(buff.Bytes()) // resets internal state
	b, err := bufio.NewReader(buff).Peek(2)
	if err != nil {
		return err
	}

	ID := session.ID ^ uint32(b[0]) ^ uint32(b[1])

	if err = binary.Write(buff, binary.BigEndian, &ID); err != nil {
		return nil
	}

	return nil
}

func (session *Session) fragmentChunks(chnks *list.List) *list.List {
	l := uint16(0)
	for c := chnks.Front(); c != nil; c = c.Next() {
		l += c.Value.(Chunk).Len()
	}

	if l <= session.mtu {
		return chnks
	}

	fragmentSlice := make([]byte, l)
	fragmentBuff := bytes.NewBuffer(fragmentSlice)
	fragmentsNum := uint16(l / session.mtu)
	if l%session.mtu > 0 {
		fragmentsNum++
	}
	for c := chnks.Front(); c != nil; c = c.Next() {
		c.Value.(Chunk).WriteTo(fragmentBuff)
	}

	pcktID := atomic.AddUint32(&session.pcktCounter, 1)
	fragmentsList := list.New()
	for i := uint16(0); i < fragmentsNum; i++ {
		fragment := &chunks.FragmentChunk{
			MoreFragments: i == fragmentsNum-1,
			PacketID:      vlu.Vlu(pcktID),
			FragmentNum:   vlu.Vlu(i),
		}

		if i == fragmentsNum-1 {
			fragment.Fragment = fragmentSlice[i*session.mtu:]
		} else {
			fragment.Fragment = fragmentSlice[i*session.mtu : (i+1)*session.mtu]
		}
		fragmentsList.PushBack(fragment)
	}
	return fragmentsList
}

// Performs defragmentation of FragmentChunk chunks.
func (session *Session) defragmentChunks(chnks *list.List) (*list.List, error) {

	var defragmentBuff = bytes.NewBuffer(make([]byte, packetMtu))
	defragmentBuff.Reset()

	pckt := &Packet{}
	pckt.writeHeaderTo(defragmentBuff)

	for chunk := chnks.Front(); chunk != nil; chunk.Next() {
		_, err := defragmentBuff.Write(chunk.Value.(*chunks.FragmentChunk).Fragment)

		if err != nil {
			return nil, err
		}
	}

	l, err := session.ReadPacket(defragmentBuff)

	return l.Chunks, err
}

// WritePacket Writes packet into the byte buffer
func (session *Session) WritePacket(pckt Packet, buff *bytes.Buffer) error {

	err := binary.Write(buff, binary.BigEndian, uint32(0))
	if err != nil {
		return err
	}

	if session.HasChecksums { // todo: error handling
		binary.Write(buff, binary.BigEndian, uint16(0))
	}

	if pckt.Len() > uint32(session.mtu) {
		pckt.Chunks, err = pckt.doFragmentation(session.mtu, &session.pcktCounter)
		if err != nil {
			return err
		}
	}

	err = pckt.writeHeaderTo(buff)
	if err != nil {
		return err
	}

	err = pckt.writeChunksTo(buff)
	if err != nil {
		return err
	}

	err = pckt.writePaddingTo(buff)
	if err != nil {
		return err
	}

	//pckt.Chunks = session.fragmentChunks(pckt.Chunks)

	//pckt.writeTo(buff)
	/*
		for c := pckt.Chunks.Front(); c != nil; c = c.Next() {
			err := pckt.writeChunkTo(c.Value.(Chunk), buff)
			if err != nil {
				return err
			}
		}*/

	if session.HasChecksums { // todo: ask & reimplement
		peeker := bufio.NewReader(buff) // default bufio buffer size is 4096
		data, err := peeker.Peek(peeker.Buffered())
		if err != nil {
			return err
		}

		if err = binary.Write(bytes.NewBuffer(buff.Bytes()),
			binary.BigEndian,
			ip.Checksum(data)); err != nil {
			return err
		}
	}

	// TODO: implement encryption as separate operation ... maybe
	/*
		err = session.encryptBuffer(buff)
		if err != nil {
			return err
		}

		err = session.writeID(buff)
	*/

	return nil
}

// ReadPacket reads packet from the byte buffer
func (session *Session) ReadPacket(buff *bytes.Buffer) (*Packet, error) {
	pckt := &Packet{}

	var err error
	//err := session.readID(buff)
	//if err != nil {
	//	return nil, err
	//}
	//err = session.decryptBuffer(buff)
	//if err != nil {
	//	return nil, err
	//}

	checksum := uint16(0)
	if session.HasChecksums {
		err := binary.Read(buff, binary.BigEndian, &checksum)
		if err != nil {
			return pckt, err
		}
	}

	if err = pckt.readFrom(buff); err != nil {
		return pckt, err
	}

	if pckt.Mode < ResponderMode {
		return pckt, errors.New("Only Responder and Startup packet modes are allowed")
	}

	if err = pckt.readChunks(buff); err != nil {
		return pckt, err
	}

	for node := pckt.Chunks.Front(); node != nil; node = node.Next() {
		switch node.Value.(Chunk).Type() {
		case chunks.FragmentChunkType:
			chunk := node.Value.(*chunks.FragmentChunk)
			if len(chunk.Fragment) < 1 {
				continue
			}
			if session.fragmented[chunk.PacketID] == nil { // todo: limit concurrent buffers for reassembling
				session.fragmented[chunk.PacketID] = CreateNewFragmentsBuffer(chunk)
			}
			buff := session.fragmented[chunk.PacketID] // todo: later RWMutex will be needed for concurrent access
			buff.Add(chunk)
			switch {
			case buff.Size > maxFragmentsSize:
				continue // todo: delete from fragmented map or not?
			case buff.Len() > maxFragments:
				continue
			case !buff.IsComplete():
				continue
			}
			// todo: reassemble packet, check it and return
		}
	}
	return pckt, nil
}

func (s *Session) reassemblePacket(f *fragmentsBuffer) (*Packet, error) {
	data := new(bytes.Buffer)
	for _, chnk := range f.buff {
		_, err := data.Write(chnk.Fragment)
		if err != nil {
			return nil, err
		}
	}
	return s.ReadPacket(data)
}
