package main

import "net"

type Target struct {
	IP   net.IP
	Port uint16
}
