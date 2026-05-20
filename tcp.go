package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// BuildSYN creates a raw IP and TCP SYN packet.
func BuildSYN(srcIP, dstIP net.IP, srcPort, dstPort uint16) ([]byte, error) {
	ipLayer := &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TOS:      0,
		Id:       54321,
		Flags:    layers.IPv4DontFragment,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(dstPort),
		SYN:     true,
		Window:  65535,
		Seq:     1105024978,
	}

	// Link the network layer for accurate checksum calculation
	if err := tcpLayer.SetNetworkLayerForChecksum(ipLayer); err != nil {
		return nil, err
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}

	err := gopacket.SerializeLayers(buf, opts, ipLayer, tcpLayer)
	return buf.Bytes(), err
}
