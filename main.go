package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	dnslookupService "dnslookup/service"
)

func main() {
	modeArg := ""
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "scan", "dnslookup", "all":
			modeArg = args[0]
			args = args[1:]
		}
	}

	config, enableDNS, dnsConfigPath, mode, err := parseConfig(modeArg, args)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	if enableDNS {
		dnsService, err := dnslookupService.New(dnsConfigPath)
		if err != nil {
			log.Fatalf("Failed to initialize DNS lookup service: %v", err)
		}
		dnsService.Start()
		log.Printf("DNS lookup service started on %s", dnsService.Config.Address)
	}

	if mode == "dnslookup" {
		select {}
	}

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

	openTargets := make(chan Target, 1000)

	// 1. Start stateful L7 workers
	WorkerPool(ctx, config.workers, openTargets)

	// 2. Start asynchronous pcap listener
	go func() {
		if err := ReceiverLoop(ctx, config.iface, config.listenPort, openTargets); err != nil {
			log.Fatalf("Receiver failed: %v", err)
		}
	}()

	// 3. Kick off raw SYN injection
	log.Printf("Using interface %s with source IP %s", config.iface, config.srcIP)
	if err := SenderLoop(ctx, config.srcIP, config.targets, config.listenPort); err != nil {
		log.Fatalf("Sender failed: %v", err)
	}

	<-ctx.Done()
}

type config struct {
	iface      string
	srcIP      net.IP
	listenPort uint16
	targets    []Target
	workers    int
}

func parseConfig(defaultMode string, args []string) (config, bool, string, string, error) {
	var (
		ifaceFlag     string
		srcIPFlag     string
		targetsFlag   string
		listenPort    uint
		workers       int
		dnsLookup     bool
		dnsConfigPath string
		mode          string
	)

	fs := flag.NewFlagSet("synkinetik", flag.ContinueOnError)
	fs.StringVar(&ifaceFlag, "iface", "", "network interface to capture on; auto-detected when omitted")
	fs.StringVar(&srcIPFlag, "src-ip", "", "local IPv4 source address; auto-detected from -iface when omitted")
	fs.StringVar(&targetsFlag, "targets", "8.8.8.8:53", "comma-separated IPv4 targets as host:port")
	fs.UintVar(&listenPort, "listen-port", 44321, "local TCP source port to use for SYN packets")
	fs.IntVar(&workers, "workers", 50, "number of TCP verification workers")
	fs.BoolVar(&dnsLookup, "dns-lookup", false, "start embedded DNS lookup service")
	fs.StringVar(&dnsConfigPath, "dns-lookup-config", "dnslookup/lookup.conf", "path to DNS lookup config file")
	fs.StringVar(&mode, "mode", "scan", "operation mode: scan, dnslookup, all")

	if defaultMode == "" {
		defaultMode = "scan"
	}
	mode = defaultMode

	if err := fs.Parse(args); err != nil {
		return config{}, false, "", "", err
	}

	if listenPort == 0 || listenPort > 65535 {
		return config{}, false, "", mode, fmt.Errorf("-listen-port must be between 1 and 65535")
	}
	if workers < 1 {
		return config{}, false, "", mode, fmt.Errorf("-workers must be at least 1")
	}

	targets, err := parseTargets(targetsFlag)
	if err != nil {
		return config{}, false, "", mode, err
	}

	iface, srcIP, err := resolveInterfaceAndSourceIP(ifaceFlag, srcIPFlag)
	if err != nil {
		return config{}, false, "", mode, err
	}

	if mode != "scan" && mode != "dnslookup" && mode != "all" {
		return config{}, false, "", mode, fmt.Errorf("invalid mode %q: must be scan, dnslookup, or all", mode)
	}

	if mode == "dnslookup" || mode == "all" {
		dnsLookup = true
	}

	return config{
		iface:      iface,
		srcIP:      srcIP,
		listenPort: uint16(listenPort),
		targets:    targets,
		workers:    workers,
	}, dnsLookup, dnsConfigPath, mode, nil
}

func parseTargets(raw string) ([]Target, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("-targets cannot be empty")
	}

	parts := strings.Split(raw, ",")
	targets := make([]Target, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		host, portRaw, err := net.SplitHostPort(part)
		if err != nil {
			return nil, fmt.Errorf("target %q must be in host:port format: %w", part, err)
		}

		ip := net.ParseIP(host).To4()
		if ip == nil {
			return nil, fmt.Errorf("target %q must use an IPv4 address", part)
		}

		port, err := strconv.ParseUint(portRaw, 10, 16)
		if err != nil || port == 0 {
			return nil, fmt.Errorf("target %q has invalid port %q", part, portRaw)
		}

		targets = append(targets, Target{IP: ip, Port: uint16(port)})
	}

	return targets, nil
}

func resolveInterfaceAndSourceIP(ifaceName, srcIPRaw string) (string, net.IP, error) {
	srcIP := net.ParseIP(srcIPRaw).To4()
	if srcIPRaw != "" && srcIP == nil {
		return "", nil, fmt.Errorf("-src-ip must be an IPv4 address")
	}

	if ifaceName != "" {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return "", nil, fmt.Errorf("interface %q not found: %w", ifaceName, err)
		}
		if srcIP != nil {
			return iface.Name, srcIP, nil
		}
		ip, err := firstIPv4(iface)
		if err != nil {
			return "", nil, err
		}
		return iface.Name, ip, nil
	}

	if srcIP != nil {
		iface, err := interfaceForIP(srcIP)
		if err != nil {
			return "", nil, err
		}
		return iface.Name, srcIP, nil
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, fmt.Errorf("failed to list network interfaces: %w", err)
	}
	for i := range ifaces {
		if ifaces[i].Flags&net.FlagUp == 0 || ifaces[i].Flags&net.FlagLoopback != 0 {
			continue
		}
		ip, err := firstIPv4(&ifaces[i])
		if err == nil {
			return ifaces[i].Name, ip, nil
		}
	}

	return "", nil, fmt.Errorf("no active non-loopback IPv4 interface found; pass -iface and -src-ip explicitly")
}

func interfaceForIP(ip net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list network interfaces: %w", err)
	}

	for i := range ifaces {
		addrs, err := ifaces[i].Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if localIP := ipNet.IP.To4(); localIP != nil && localIP.Equal(ip) {
					return &ifaces[i], nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no interface has source IP %s; pass -iface explicitly", ip)
}

func firstIPv4(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses for %s: %w", iface.Name, err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			ip := ipNet.IP.To4()
			if ip != nil && !ip.IsLoopback() {
				return ip, nil
			}
		}
	}
	return nil, fmt.Errorf("interface %s has no non-loopback IPv4 address", iface.Name)
}
