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

package protocol

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"github.com/rtmfpew/rtmfpew/protocol/chunks"
	"github.com/rtmfpew/rtmfpew/protocol/crypto"
	"github.com/rtmfpew/rtmfpew/protocol/io"
	"github.com/rtmfpew/rtmfpew/protocol/net"
	"sync/atomic"
)

const defaultMtu = 768

const (
	// NormalSession is a client-server
	NormalSession = 1 << (iota)
	// ForwardedSession is a NAT traversed session
	ForwardedSession
	// RedirectedSession is a peer redirected session
	RedirectedSession
	// RendezvousSession is a P2P session
	RendezvousSession
)

// Session stores current connection state
type Session struct {
	ID uint32

	initiatorAddr *net.PeerAddress
	responderAddr *net.PeerAddress

	profile *crypto.Profile
	key     []byte

	pcktCounter uint32
	mtu         uint16

	hasChecksums bool
}

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

func (session *Session) writeID(buff *bytes.Buffer, b []byte) error {
	ID := session.ID ^ uint32(b[0]) ^ uint32(b[1])
	err := binary.Write(buff, binary.BigEndian, &ID)
	if err != nil {
		return nil
	}

	return nil
}

func (session *Session) fragmentChunks(chnks *list.List) *list.List {
	fragmentSlice := make([]byte, defaultMtu)
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
				PacketID:      io.Vlu(pcktID),
				FragmentNum:   io.Vlu(i),
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

var defragmentBuff = bytes.NewBuffer(make([]byte, defaultMtu))

func (session *Session) defragmentChunks(chnks []chunks.FragmentChunk) (*list.List, error) {
	defragmentBuff.Reset()

	pckt := &Packet{}
	pckt.writeTo(defragmentBuff)

	for i := 0; i < len(chnks); i++ {
		_, err := defragmentBuff.Write(chnks[i].Fragment)

		if err != nil {
			return nil, err
		}
	}

	l, err := session.ReadPacket(defragmentBuff)

	return l.Chunks, err
}

// WritePacket Writes packet into the byte buffer
func (session *Session) WritePacket(pckt Packet, buff *bytes.Buffer) error {

	if session.hasChecksums {
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

	if session.hasChecksums {
		if err = binary.Write(buff,
			binary.BigEndian,
			io.Checksum(data)); err != nil {
			return err
		}
	}

	return nil
}

// ReadPacket reads packet from the byte buffer
func (session *Session) ReadPacket(buff *bytes.Buffer) (*Packet, error) {
	pckt := &Packet{}

	pckt.Chunks = list.New()

	checksum := uint16(0)
	if session.hasChecksums {
		err := binary.Read(buff, binary.BigEndian, &checksum)
		if err != nil {
			return pckt, err
		}
	}

	err := pckt.readFrom(buff)
	if err != nil {
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
		case chunks.FragmentChunkType:
			c := &chunks.FragmentChunk{}
			if err = c.ReadFrom(buff); err != nil {
				return pckt, err
			}
			datalen += c.Len()
			pckt.Chunks.PushBack(c)
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

	if session.hasChecksums {
		peeker := bufio.NewReader(buff)
		data, err := peeker.Peek(peeker.Buffered())
		if err != nil {
			return pckt, err
		}

		calcedChecksum := io.Checksum(data)
		if calcedChecksum != checksum {
			return pckt, errors.New("Wrong packet checksum")
		}
	}

	return pckt, nil
}
