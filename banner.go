package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

// WorkerPool consumes open targets and grabs service banners concurrently.
func WorkerPool(ctx context.Context, concurrency int, targets <-chan Target) {
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-targets:
					if !ok {
						return // Channel closed
					}
					grabBanner(t)
				}
			}
		}(i)
	}
}

func grabBanner(t Target) {
	addr := fmt.Sprintf("%s:%d", t.IP, t.Port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	// Stub: In a full implementation, read the banner data here and regex match.
	log.Printf("[+] Open Target Found & Verified: %s", addr)
}
