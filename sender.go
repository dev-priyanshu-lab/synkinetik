package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"syscall"
)

// SenderLoop runs asynchronously, crafting and injecting raw packets.
func SenderLoop(ctx context.Context, srcIP net.IP, targets []Target, srcPort uint16) error {
	// Open a raw socket (macOS / BSD compatible)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return fmt.Errorf("failed to open raw socket: %w", err)
	}
	defer syscall.Close(fd)

	// IP_HDRINCL tells the kernel we are providing the IP header
	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		return fmt.Errorf("failed to set IP_HDRINCL: %w", err)
	}

	log.Printf("Sender loop started to %d targets...", len(targets))

	for _, t := range targets {
		select {
		case <-ctx.Done():
			log.Println("Sender loop terminating...")
			return nil
		default:
			rawPacket, err := BuildSYN(srcIP, t.IP, srcPort, t.Port)
			if err != nil {
				log.Printf("Error building packet for %s: %v", t.IP, err)
				continue
			}

			var addr [4]byte
			copy(addr[:], t.IP.To4())

			sockAddr := &syscall.SockaddrInet4{
				Port: 0, // Port is in the raw packet
				Addr: addr,
			}

			if err := syscall.Sendto(fd, rawPacket, 0, sockAddr); err != nil {
				log.Printf("Failed to send packet to %s: %v", t.IP, err)
			}
		}
	}
	return nil
}
