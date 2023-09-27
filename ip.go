package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"net"
	"sync"
	"time"
)

// SystemPortsRange is the range of system ports to scan
const SystemPortsRange = 1024

// Port is a struct for a port
type Port struct {
	ID         int    `json:"id"`
	PortNumber string `json:"text"`
}

// PortTree is a tree of IPs and their pinged ports
type PortTree struct {
	Ports []Port `json:"children"`
	ID    string `json:"id"`
	IP    string `json:"text"`
}

// pingHost pings a single host
func pingHost(ip string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("Error creating pinger for %s: %v\n", ip, err)
		return
	}
	pinger.Count = 1 // Send a single ICMP packet
	pinger.Timeout = time.Second

	err = pinger.Run()
	if err != nil {
		return
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		results <- ip
	}
}

// pingIPRange pings all hosts in a slice of IPs
func (a *App) pingIPRange(ips []string) []string {
	var wg sync.WaitGroup
	results := make(chan string, len(ips))

	for _, ip := range ips {
		wg.Add(1)
		go pingHost(ip, &wg, results)
	}
	wg.Wait()
	close(results)

	var activeIPs []string
	for ip := range results {
		activeIPs = append(activeIPs, ip)
	}
	fmt.Printf("Active IPs: %v\n", activeIPs)
	return activeIPs
}

// inc increments the IP address by 1
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// PingPort pings a port on a host
func (a *App) PingPort(ip string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	fmt.Printf("Pinging %s:%d\n", ip, port)
	if err != nil {
		fmt.Printf("Error pinging %s:%d: %v\n", ip, port, err)
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

// generateIPsForCIDR takes a CIDR and generates all IP addresses in that range
func (a *App) generateIPsForCIDR(cidr string) ([]string, error) {
	fmt.Printf("Generating range for %s\n", cidr)
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// Remove network and broadcast addresses
	return ips[1 : len(ips)-1], nil
}

// PingIPsByCIDR pings all hosts in a CIDR range
func (a *App) PingIPsByCIDR(cidr string) []string {
	ips, err := a.generateIPsForCIDR(cidr)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	fmt.Printf("Pinging %d hosts\n", len(ips))
	return a.pingIPRange(ips)
}

// pingPortsForSingleIP pings all system ports for a given IP
func (a *App) pingPortsForSingleIP(ip string, results chan<- PortTree, ipSem, portSem chan struct{}) {
	defer func() { <-ipSem }() // Release IP semaphore

	tree := PortTree{
		ID:    ip,
		IP:    ip,
		Ports: []Port{},
	}

	var wg sync.WaitGroup

	portResults := make(chan Port, SystemPortsRange)

	// Iterate through system ports
	for port := 1; port <= SystemPortsRange; port++ {
		portSem <- struct{}{} // Acquire port semaphore
		wg.Add(1)

		go func(port int) {
			defer wg.Done()
			defer func() { <-portSem }() // Release port semaphore

			if a.PingPort(ip, port) {
				portResults <- Port{
					ID:         port,
					PortNumber: fmt.Sprintf("%d", port),
				}
			}
		}(port)
	}

	// Wait for all port checks to complete for this IP
	go func() {
		wg.Wait()
		close(portResults)
	}()

	// Collect results
	for port := range portResults {
		tree.Ports = append(tree.Ports, port)
	}

	results <- tree
}

// getTreeForIPsRange pings all system ports for a slice of IPs
func (a *App) getTreeForIPsRange(ips []string, ipSem, portSem chan struct{}) []PortTree {
	var trees []PortTree

	results := make(chan PortTree, len(ips))

	for _, ip := range ips {
		ipSem <- struct{}{} // Acquire IP semaphore
		go a.pingPortsForSingleIP(ip, results, ipSem, portSem)
	}

	for i := 0; i < len(ips); i++ {
		trees = append(trees, <-results)
	}

	return trees
}

// GeneratePortTreesByCIDR generates a tree of IPs and their pinged ports for a given CIDR
func (a *App) GeneratePortTreesByCIDR(cidr string) []PortTree {
	ips := a.PingIPsByCIDR(cidr)
	fmt.Printf("Pinged %d hosts\n", len(ips))
	// Semaphore to limit concurrent goroutines for IP scanning
	var ipSem = make(chan struct{}, len(ips)/2) // e.g., len(ips)/2 concurrent IP scans

	// Semaphore to limit concurrent goroutines for Port scanning
	var portSem = make(chan struct{}, len(ips)) // e.g., len(ips) concurrent port scans
	return a.getTreeForIPsRange(ips, ipSem, portSem)
}

// GeneratePortTreeForIP generates a tree of IPs and their pinged ports for a given IP
func (a *App) GeneratePortTreeForIP(ip string) []PortTree {
	// Semaphore to limit concurrent goroutines for IP scanning
	var ipSem = make(chan struct{}, 1) // e.g., len(ips)/2 concurrent IP scans

	// Semaphore to limit concurrent goroutines for Port scanning
	var portSem = make(chan struct{}, 1) // e.g., len(ips) concurrent port scans
	trees := a.getTreeForIPsRange([]string{ip}, ipSem, portSem)
	fmt.Println(trees)
	return trees
}

// AutoDetect detects the local machine's IP and its subnet mask, then returns a CIDR
func (a *App) AutoDetect() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if ipv4 := ip.To4(); ipv4 != nil && !ipv4.IsLoopback() {
				return addr.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable IP address found")
}
