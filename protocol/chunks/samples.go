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

package chunks

import (
	"net"

	"github.com/rtmfpew/amfy/vlu"
	"github.com/rtmfpew/rtmfpew/protocol/connection"
)

// BufferProbeChunkSample returns buffer probe sample chunk
func BufferProbeChunkSample() *BufferProbeChunk {
	return &BufferProbeChunk{
		FlowID: vlu.Vlu(16),
	}
}

// DataAcknowledgementBitmapChunkSample returns data acknowledgement bitmap sample chunk
func DataAcknowledgementBitmapChunkSample() *DataAcknowledgementBitmapChunk {
	data := [...]byte{
		0x16, 0x11, 0x1A, 0x3A,
		0x1B, 0x5C, 0xAC, 0x7E,
	}

	return &DataAcknowledgementBitmapChunk{
		FlowID:                vlu.Vlu(19872),
		BufferBlocksAvailable: vlu.Vlu(1141),
		CumulativeAck:         vlu.Vlu(391),
		Acknowledgement:       data[:],
	}
}

// DataAcknowledgementRangesChunkSample returns data acknowledgement range sample chunk
func DataAcknowledgementRangesChunkSample() *DataAcknowledgementRangesChunk {
	ranges := [...]DataAcknowledgementRange{
		DataAcknowledgementRange{
			HolesMinusOne:    vlu.Vlu(141),
			ReceivedMinusOne: vlu.Vlu(1492351),
		}, DataAcknowledgementRange{
			HolesMinusOne:    vlu.Vlu(1112),
			ReceivedMinusOne: vlu.Vlu(119),
		}, DataAcknowledgementRange{
			HolesMinusOne:    vlu.Vlu(10131),
			ReceivedMinusOne: vlu.Vlu(109),
		}, DataAcknowledgementRange{
			HolesMinusOne:    vlu.Vlu(151),
			ReceivedMinusOne: vlu.Vlu(1191),
		},
	}

	return &DataAcknowledgementRangesChunk{
		FlowID:                vlu.Vlu(1831),
		CumulativeAck:         vlu.Vlu(119),
		BufferBlocksAvailable: vlu.Vlu(284),
		Ranges:                ranges[:],
	}
}

// FlowExceptionReportChunkSample returns flow exception report sample chunk
func FlowExceptionReportChunkSample() *FlowExceptionReportChunk {
	return &FlowExceptionReportChunk{
		FlowID:    vlu.Vlu(182),
		Exception: vlu.Vlu(161),
	}
}

// ForwardedHelloChunkSample returns forwarded hello sample chunk
func ForwardedHelloChunkSample() *ForwardedHelloChunk {
	addr, _ := net.ResolveUDPAddr("udp", "192.168.1.1:1935")

	epd := [...]byte{0x91, 0xF1, 0xAA, 0xBC, 0xAD}
	tag := [...]byte{0x1A, 0xB2, 0xBA, 0xDC, 0xED}

	return &ForwardedHelloChunk{
		Epd:          epd[:],
		Tag:          tag[:],
		ReplyAddress: *connection.PeerAddressFrom(addr),
	}
}

// FragmentChunkSample returns fragment sample chunk
func FragmentChunkSample() *FragmentChunk {

	frag := [...]byte{0x12, 0x9A, 0x1A, 0xFF}
	return &FragmentChunk{
		MoreFragments: true,
		PacketID:      vlu.Vlu(123),
		FragmentNum:   vlu.Vlu(231),
		Fragment:      frag[:],
	}
}

// HelloCookieChangeChunkSample returns hello cookie change sample chunk
func HelloCookieChangeChunkSample() *HelloCookieChangeChunk {

	oldC := [...]byte{0x11, 0xBA, 0x2A, 0xEF, 0xA1}
	newC := [...]byte{0x12, 0x9A, 0x1A, 0xFD, 0x91}
	return &HelloCookieChangeChunk{
		OldCookie: oldC[:],
		NewCookie: newC[:],
	}
}

// InitiatorHelloChunkSample returns hello cookie change sample chunk
func InitiatorHelloChunkSample() *InitiatorHelloChunk {

	epd := [...]byte{0x91, 0xF1, 0xAA, 0xBC, 0xAD}
	tag := [...]byte{0x1A, 0xB2, 0xBA, 0xDC, 0xED}

	return &InitiatorHelloChunk{
		Epd: epd[:],
		Tag: tag[:],
	}
}

// InitiatorInitialKeyingChunkSample returns initiator initial keying sample chunk
func InitiatorInitialKeyingChunkSample() *InitiatorInitialKeyingChunk {
	data := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xA1}

	return &InitiatorInitialKeyingChunk{
		CookieEcho:                   data[:],
		InitiatorCertificate:         data[:],
		InitiatorSessionID:           112,
		SessionKeyInitiatorComponent: data[:],
		Signature:                    data[:],
	}
}

// PingReplyChunkSample returns ping reply sample chunk
func PingReplyChunkSample() *PingReplyChunk {
	m := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xA1}

	return &PingReplyChunk{
		MessageEcho: m[:],
	}
}

// PingChunkSample returns ping reply sample chunk
func PingChunkSample() *PingChunk {
	msg := [...]byte{0x2B, 0xC3, 0xB1, 0x5C, 0xED, 0xA1}

	return &PingChunk{
		Message: msg[:],
	}
}

// ResponderHelloChunkSample returns responder hello sample chunk
func ResponderHelloChunkSample() *ResponderHelloChunk {
	data := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xA2, 0x18, 0xA1, 0xB2, 0x6F}

	return &ResponderHelloChunk{
		Cookie:               data[:],
		ResponderCertificate: data[:],
		TagEcho:              data[:],
	}
}

// ResponderInitialKeyingChunkSample returns responder initial keying sample chunk
func ResponderInitialKeyingChunkSample() *ResponderInitialKeyingChunk {
	data := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xA2, 0x18, 0xA1, 0xB2}

	return &ResponderInitialKeyingChunk{
		ResponderSessionID:           23812,
		SessionKeyResponderComponent: data[:],
		Signature:                    data[:],
	}
}

// ResponderRedirectChunkSample returns responder initial keying sample chunk
func ResponderRedirectChunkSample() *ResponderRedirectChunk {

	data := [...]byte{0x2A, 0xC3, 0xB1, 0x5C}
	addresses := [...]connection.PeerAddress{
		connection.PeerAddress{
			IP:     data[:],
			Port:   2913,
			Origin: connection.RemoteOrigin,
		}, connection.PeerAddress{
			IP:     data[:],
			Port:   2911,
			Origin: connection.ProxyOrigin,
		}, connection.PeerAddress{
			IP:     data[:],
			Port:   2912,
			Origin: connection.LocalOrigin,
		},
	}

	return &ResponderRedirectChunk{
		RedirectDestination: addresses[:],
		TagEcho:             data[:],
	}
}

// SessionCloseAcknowledgementSample returns session close acknowledgement sample chunk
func SessionCloseAcknowledgementSample() *SessionCloseAcknowledgement {
	return &SessionCloseAcknowledgement{}
}

// SessionCloseRequestChunkSample returns session close request sample chunk
func SessionCloseRequestChunkSample() *SessionCloseRequestChunk {
	return &SessionCloseRequestChunk{}
}

// UserDataOptionSample returns user data option sample chunk
func UserDataOptionSample() *UserDataOption {

	data := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xA2, 0x18, 0xA1, 0xB2}
	return &UserDataOption{
		OptionType: vlu.Vlu(912),
		Value:      data[:],
	}
}

// UserDataChunkSample returns user data option sample chunk
func UserDataChunkSample() *UserDataChunk {

	t := [...]byte{0x2A, 0xC3, 0xB1, 0x5C, 0xED, 0xEC}
	opts := [...]UserDataOption{
		UserDataOption{
			OptionType: 1,
			Value:      t[:],
		}, UserDataOption{
			OptionType: 2,
			Value:      t[:],
		},
	}

	return &UserDataChunk{
		Abandon:         true,
		Final:           false,
		FlowID:          vlu.Vlu(126),
		FragmentControl: MiddleFragmentControl,
		FsnOffset:       vlu.Vlu(12),
		Options:         opts[:],
		OptionsPresent:  true,
		SequenceNumber:  vlu.Vlu(1256),
		UserData:        t[:],
	}
}
