// Copyright 2016 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package layers

import (
	"reflect"
	"testing"

	"github.com/helshabini/gopacket"
)

// testPacketMPLS
// Ethernet II, Src: cc:15:14:64:00:00 (cc:15:14:64:00:00), Dst: cc:13:14:64:00:01 (cc:13:14:64:00:01)
// MultiProtocol Label Switching Header, Label: 17, Exp: 0, S: 0, TTL: 254
// MultiProtocol Label Switching Header, Label: 19, Exp: 0, S: 1, TTL: 254
// Internet Protocol Version 4, Src: 12.0.0.1, Dst: 2.2.2.2
// Internet Control Message Protocol
// 0000   cc 13 14 64 00 01 cc 15 14 64 00 00 88 47 00 01  ...d.....d...G..
// 0010   10 fe 00 01 31 fe 45 00 00 64 00 39 00 00 fe 01  ....1.E..d.9....
// 0020   ac 5b 0c 00 00 01 02 02 02 02 08 00 3a 6b 00 0b  .[..........:k..
// 0030   00 02 00 00 00 00 00 3e 43 94 ab cd ab cd ab cd  .......>C.......
// 0040   ab cd ab cd ab cd ab cd ab cd ab cd ab cd ab cd  ................
// 0050   ab cd ab cd ab cd ab cd ab cd ab cd ab cd ab cd  ................
// 0060   ab cd ab cd ab cd ab cd ab cd ab cd ab cd ab cd  ................
// 0070   ab cd ab cd ab cd ab cd ab cd                    ..........

var testPacketMPLS = []byte{
	0xcc, 0x13, 0x14, 0x64, 0x00, 0x01, 0xcc, 0x15, 0x14, 0x64, 0x00, 0x00, 0x88, 0x47, 0x00, 0x01,
	0x10, 0xfe, 0x00, 0x01, 0x31, 0xfe, 0x45, 0x00, 0x00, 0x64, 0x00, 0x39, 0x00, 0x00, 0xfe, 0x01,
	0xac, 0x5b, 0x0c, 0x00, 0x00, 0x01, 0x02, 0x02, 0x02, 0x02, 0x08, 0x00, 0x3a, 0x6b, 0x00, 0x0b,
	0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x43, 0x94, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
	0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
	0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
	0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
	0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
}

func TestPacketMPLS(t *testing.T) {
	p := gopacket.NewPacket(testPacketMPLS, LinkTypeEthernet, gopacket.Default)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeMPLS, LayerTypeMPLS, LayerTypeIPv4, LayerTypeICMPv4, gopacket.LayerTypePayload}, t)
	if got, ok := p.Layers()[1].(*MPLS); ok {
		want := &MPLS{
			BaseLayer: BaseLayer{
				Contents: []byte{0x00, 0x01, 0x10, 0xfe},
				Payload: []byte{0x00, 0x01, 0x31, 0xfe, 0x45, 0x00, 0x00, 0x64, 0x00, 0x39, 0x00, 0x00, 0xfe, 0x01,
					0xac, 0x5b, 0x0c, 0x00, 0x00, 0x01, 0x02, 0x02, 0x02, 0x02, 0x08, 0x00, 0x3a, 0x6b, 0x00, 0x0b,
					0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x43, 0x94, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd},
			},
			Label:        17,
			TrafficClass: 0,
			StackBottom:  false,
			TTL:          254,
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("MPLS layer 1 mismatch, \nwant %#v\ngot %#v\n", want, got)
		}
	}
	if got, ok := p.Layers()[2].(*MPLS); ok {
		want := &MPLS{
			BaseLayer: BaseLayer{
				Contents: []byte{0x00, 0x01, 0x31, 0xfe},
				Payload: []byte{0x45, 0x00, 0x00, 0x64, 0x00, 0x39, 0x00, 0x00, 0xfe, 0x01,
					0xac, 0x5b, 0x0c, 0x00, 0x00, 0x01, 0x02, 0x02, 0x02, 0x02, 0x08, 0x00, 0x3a, 0x6b, 0x00, 0x0b,
					0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x43, 0x94, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
					0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd},
			},
			Label:        19,
			TrafficClass: 0,
			StackBottom:  true,
			TTL:          254,
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("MPLS layer 2 mismatch, \nwant %#v\ngot %#v\n", want, got)
		}
	}
}

func BenchmarkDecodePacketMPLS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gopacket.NewPacket(testPacketMPLS, LinkTypeEthernet, gopacket.NoCopy)
	}
}
