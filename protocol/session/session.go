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
	"github.com/rtmfpew/rtmfpew/config"
	"github.com/rtmfpew/rtmfpew/protocol/chunks"
	"github.com/rtmfpew/rtmfpew/protocol/crypto"
	"github.com/rtmfpew/rtmfpew/protocol/ip"
	"github.com/rtmfpew/rtmfpew/protocol/net"
	"github.com/rtmfpew/rtmfpew/protocol/vlu"
	"sync/atomic"
)

var (
	packetMtu           = config.Mtu()
	maxFragmentationGap = config.MaxFragmentationGap()
	maxFragments        = config.MaxFragments()
)

const (
	// NormalSession is a client-server
	NormalSession = 1 << (iota)
	// ForwardedSession is a NAT traversed session
	ForwardedSession
	// RedirectedSession is a peer redirected session
	RedirectedSession
	// RendezvousSession is a P2P session
	RendezvousSession
	// HandshakeSession is an initial session handshake
	HandshakeSession
)

// Session stores current connection state
type Session struct {
	ID uint32

	InitiatorAddr *net.PeerAddress
	ResponderAddr *net.PeerAddress

	profile crypto.Profile

	pcktCounter uint32
	mtu         uint16

	HasChecksums bool
	Established  bool

	fragments *list.List
}

// NewSessionWith creates new session with custom profile
func NewSessionWith(profile crypto.Profile) *Session {
	session := &Session{
		profile:      profile,
		pcktCounter:  0,
		mtu:          uint16(packetMtu),
		HasChecksums: false,
		Established:  false,
	}

	return session
}

// NewSession creates new session with default profile
func NewSession() *Session {
	profile := &crypto.DefaultProfile{}
	profile.InitDefault()

	return NewSessionWith(profile)
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
func (session *Session) readID(buff *bytes.Buffer) error {
	err := binary.Read(buff, binary.BigEndian, &session.ID)
	if err != nil {
		return err
	}

	peeker := bufio.NewReader(buff)
	b, err := peeker.Peek(2) // Unscramble session ID
	session.ID = session.ID ^ uint32(b[0]) ^ uint32(b[1])

	return nil
}

// writeID writes session ID
func (session *Session) writeID(buff *bytes.Buffer) error {
	buff = bytes.NewBuffer(buff.Bytes()) // resets internal state
	peeker := bufio.NewReader(buff)
	b, err := peeker.Peek(2)
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
	fragmentSlice := make([]byte, session.mtu)
	fragmentBuff := bytes.NewBuffer(fragmentSlice)

	l := uint16(0)
	for c := chnks.Front(); c != nil; c = c.Next() {
		l += c.Value.(Chunk).Len()
	}

	if l > session.mtu {
		fragmentsNum := uint16(l / session.mtu)
		if l%session.mtu > 0 {
			fragmentsNum++
		}

		fragments := make([]chunks.FragmentChunk, fragmentsNum)

		for c := chnks.Front(); c != nil; c = c.Next() {
			c.Value.(Chunk).WriteTo(fragmentBuff)
		}

		pcktID := atomic.AddUint32(&session.pcktCounter, 1)

		for i := uint16(0); i < fragmentsNum; i++ {
			fragment := &chunks.FragmentChunk{
				MoreFragments: i == fragmentsNum-1,
				PacketID:      vlu.Vlu(pcktID),
				FragmentNum:   vlu.Vlu(i),
			}

			if i == fragmentsNum-1 {
				fragment.Fragment = fragmentSlice[i*session.mtu:]
				break
			}

			fragment.Fragment = fragmentSlice[i*session.mtu : (i+1)*session.mtu]
			fragments[i] = *fragment
		}

		fragmentsList := list.New()

		for _, fragment := range fragments {
			fragmentsList.PushBack(fragment)
		}

		return fragmentsList
	}

	return nil
}

var defragmentBuff = bytes.NewBuffer(make([]byte, packetMtu))

func (session *Session) defragmentChunks(chnks *list.List) (*list.List, error) {
	defragmentBuff.Reset()

	pckt := &Packet{}
	pckt.writeTo(defragmentBuff)

	for chunk := chnks.Front(); chunk != nil; chunk.Next() {
		_, err := defragmentBuff.Write(chunk.Value.(chunks.FragmentChunk).Fragment)

		if err != nil {
			return nil, err
		}
	}

	l, err := session.ReadPacket(defragmentBuff)

	return l.Chunks, err
}

func (session *Session) removeFragmentsWithPacketID(packetID vlu.Vlu) {
	for sessionFragment := session.fragments.Front(); sessionFragment != nil; sessionFragment.Next() {
		if sessionFragment.Value.(chunks.FragmentChunk).PacketID == packetID {
			session.fragments.Remove(sessionFragment)
		}
	}
}

func (session *Session) chunksToDefragment() *list.List {
	if session.fragments.Len() == 0 {
		return nil
	}

	packets := list.New()

	// Count packet IDs
	for fragment := session.fragments.Front(); fragment != nil; fragment.Next() {
		found := false
		for packet := packets.Front(); packet != nil; packet.Next() {
			if found = packet.Value.(vlu.Vlu) == fragment.Value.(chunks.FragmentChunk).PacketID; found {
				break
			}
		}

		if !found {
			packets.PushBack(fragment.Value.(chunks.FragmentChunk).PacketID)
		}
	}

	// Group packet fragments
	packetFragments := make(map[vlu.Vlu]*list.List)
	for packet := packets.Front(); packet != nil; packet.Next() {
		packetFragments[packet.Value.(vlu.Vlu)] = list.New()
		for fragment := session.fragments.Front(); fragment != nil; fragment.Next() {
			if packet.Value.(vlu.Vlu) == fragment.Value.(chunks.FragmentChunk).PacketID {
				packetFragments[packet.Value.(vlu.Vlu)].PushBack(fragment)
			}
		}

		if packetFragments[packet.Value.(vlu.Vlu)].Len() > int(maxFragments) {
			session.removeFragmentsWithPacketID(packet.Value.(vlu.Vlu))
		}
	}

	// Naive fragments sorting
	for packetID, fragments := range packetFragments {
		firstFragment := fragments.Front()
		secondFragment := fragments.Front()
		secondFragment.Next()

		if secondFragment != nil {
			permutations := 0
			for {
				if secondFragment != nil {
					// Delete duplicates
					if firstFragment.Value.(chunks.FragmentChunk).FragmentNum == secondFragment.Value.(chunks.FragmentChunk).FragmentNum {
						fragments.Remove(firstFragment)
						permutations++
					} else {
						if firstFragment.Value.(chunks.FragmentChunk).FragmentNum > secondFragment.Value.(chunks.FragmentChunk).FragmentNum {
							fragments.MoveAfter(firstFragment, secondFragment)
							permutations++
						} else {
							// Corrupted packet
							fragmentationGap := secondFragment.Value.(chunks.FragmentChunk).FragmentNum - firstFragment.Value.(chunks.FragmentChunk).FragmentNum

							if fragmentationGap != 1 && fragmentationGap >= vlu.Vlu(maxFragmentationGap) {
								session.removeFragmentsWithPacketID(packetID)
								delete(packetFragments, packetID)
							}

							return nil
						}
					}

					firstFragment.Next()
					secondFragment.Next()
				} else {
					firstFragment = fragments.Front()
					secondFragment = fragments.Front()
					secondFragment.Next()
					if permutations == 0 {
						break
					} else {
						permutations = 0
					}
				}
			}

			if !fragments.Back().Value.(chunks.FragmentChunk).MoreFragments {
				session.removeFragmentsWithPacketID(packetID)
				delete(packetFragments, packetID)
				return fragments
			}

		} else {

			// Only one fragment ?
			if !firstFragment.Value.(chunks.FragmentChunk).MoreFragments {
				session.removeFragmentsWithPacketID(packetID)
				delete(packetFragments, packetID)
				return fragments
			}
		}
	}

	return nil
}

// WritePacket Writes packet into the byte buffer
func (session *Session) WritePacket(pckt Packet, buff *bytes.Buffer) error {

	binary.Write(buff, binary.BigEndian, uint32(0))

	if session.HasChecksums {
		binary.Write(buff, binary.BigEndian, uint16(0))
	}

	pckt.Chunks = session.fragmentChunks(pckt.Chunks)

	pckt.writeTo(buff)
	for c := pckt.Chunks.Front(); c != nil; c = c.Next() {
		err := pckt.writeChunkTo(c.Value.(Chunk), buff)
		if err != nil {
			return err
		}
	}

	err := pckt.writePaddingTo(buff)
	if err != nil {
		return err
	}

	peeker := bufio.NewReader(buff)
	data, err := peeker.Peek(peeker.Buffered())
	if err != nil {
		return err
	}

	if session.HasChecksums {
		if err = binary.Write(bytes.NewBuffer(buff.Bytes()),
			binary.BigEndian,
			ip.Checksum(data)); err != nil {
			return err
		}
	}

	err = session.encryptBuffer(buff)
	if err != nil {
		return err
	}

	err = session.writeID(buff)

	return err
}

// ReadPacket reads packet from the byte buffer
func (session *Session) ReadPacket(buff *bytes.Buffer) (*Packet, error) {
	pckt := &Packet{}

	err := session.readID(buff)
	if err != nil {
		return nil, err
	}

	err = session.decryptBuffer(buff)
	if err != nil {
		return nil, err
	}

	pckt.Chunks = list.New()

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

	datalen := uint16(0)

	typ := byte(0)
	for {

		if typ, err = buff.ReadByte(); err != nil {
			return pckt, err
		}

		if typ == 0xFF || typ == 0x00 {
			break
		}

		rightType := true

		switch typ {
		case chunks.BufferProbeChunkType:
			c := &chunks.BufferProbeChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.DataAcknowledgementBitmapChunkType:
			c := &chunks.DataAcknowledgementBitmapChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.DataAcknowledgementRangesChunkType:
			c := &chunks.DataAcknowledgementRangesChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.FlowExceptionReportChunkType:
			c := &chunks.FlowExceptionReportChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.ForwardedHelloChunkType:
			c := &chunks.ForwardedHelloChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		// Store fragments
		case chunks.FragmentChunkType:
			c := &chunks.FragmentChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			session.fragments.PushBack(c)

			fragments := session.chunksToDefragment()
			if fragments == nil {
				fragments = session.chunksToDefragment()
			}

			if fragments != nil {
				chnks, err := session.defragmentChunks(fragments)
				if err != nil {
					return pckt, err
				}

				for chunk := chnks.Front(); chunk != nil; chunk.Next() {
					pckt.Chunks.PushBack(chunk.Value.(chunks.FragmentChunk))
				}
			}

			datalen += c.Len()
			break

		case chunks.HelloCookieChangeChunkType:
			c := &chunks.HelloCookieChangeChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.InitiatorHelloChunkType:
			c := &chunks.InitiatorHelloChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.InitiatorInitialKeyingChunkType:
			c := &chunks.InitiatorInitialKeyingChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.NextUserDataChunkType:
			c := &chunks.NextUserDataChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.PingReplyChunkType:
			c := &chunks.PingReplyChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.PingChunkType:
			c := &chunks.PingChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.ResponderHelloChunkType:
			c := &chunks.ResponderHelloChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.ResponderInitialKeyingChunkType:
			c := &chunks.ResponderInitialKeyingChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.ResponderRedirectChunkType:
			c := &chunks.ResponderRedirectChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.SessionCloseAcknowledgementType:
			c := &chunks.SessionCloseAcknowledgement{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.SessionCloseRequestChunkType:
			c := &chunks.SessionCloseRequestChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		case chunks.UserDataChunkType:
			c := &chunks.UserDataChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
			break

		default:
			rightType = false
			break
		}

		if !rightType {
			break
		}
	}

	if session.HasChecksums {
		peeker := bufio.NewReader(buff)
		data, err := peeker.Peek(peeker.Buffered())
		if err != nil {
			return pckt, err
		}

		calcedChecksum := ip.Checksum(data)
		if calcedChecksum != checksum {
			return pckt, errors.New("Wrong packet checksum")
		}
	}

	return pckt, nil
}
