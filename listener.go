package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ReceiverLoop captures SYN-ACKs and forwards valid targets to the results channel.
func ReceiverLoop(ctx context.Context, iface string, listenPort uint16, results chan<- Target) error {
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %w", iface, err)
	}
	defer handle.Close()

	// Strict BPF filter: only incoming SYN-ACKs to our specific source port
	bpfFilter := fmt.Sprintf("tcp and dst port %d and (tcp[tcpflags] & tcp-ack != 0)", listenPort)
	if err := handle.SetBPFFilter(bpfFilter); err != nil {
		return fmt.Errorf("failed to set BPF filter: %w", err)
	}

	log.Printf("Receiver loop listening on %s (port %d)...", iface, listenPort)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	in := packetSource.Packets()

	for {
		select {
		case <-ctx.Done():
			log.Println("Receiver loop terminating...")
			return nil
		case p := <-in:
			if p == nil {
				continue
			}
			tcpLayer := p.Layer(layers.LayerTypeTCP)
			ipLayer := p.Layer(layers.LayerTypeIPv4)

			if tcpLayer != nil && ipLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				ip, _ := ipLayer.(*layers.IPv4)
				if tcp.SYN && tcp.ACK {
					// Found an open port! Send to fingerprinter.
					results <- Target{IP: ip.SrcIP, Port: uint16(tcp.SrcPort)}
				}
			}
		}
	}
}
