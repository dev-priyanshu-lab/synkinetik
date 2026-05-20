package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup graceful shutdown boundaries
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received termination signal. Shutting down...")
		cancel()
	}()

	// Configuration stubs
	listenPort := uint16(44321)
	iface := "en0"                        // Assumed primary macOS interface
	srcIP := net.ParseIP("192.168.1.100") // Stub: Ideally resolved dynamically

	openTargets := make(chan Target, 1000)

	// 1. Start stateful L7 workers
	WorkerPool(ctx, 50, openTargets)

	// 2. Start asynchronous pcap listener
	go func() {
		if err := ReceiverLoop(ctx, iface, listenPort, openTargets); err != nil {
			log.Fatalf("Receiver failed: %v", err)
		}
	}()

	// 3. Kick off raw SYN injection
	testTargets := []Target{{IP: net.ParseIP("8.8.8.8"), Port: 53}}
	if err := SenderLoop(ctx, srcIP, testTargets, listenPort); err != nil {
		log.Fatalf("Sender failed: %v", err)
	}

	<-ctx.Done()
}
