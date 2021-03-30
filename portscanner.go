/* Copyright 2021 Roland Heimpoldinger <roland@heimpoldinger.net>

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Scan struct {
	Targets []string
	Ports   []int
	Threads int
}

var scan Scan
var portParam string
var targets []string

var top20 = []int{21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3380, 5900, 8080}

func init() {
	flag.IntVar(&scan.Threads, "t", 100, "set concurrent threads")
	flag.StringVar(&portParam, "p", "top20", "set port(s) to scan \ns: top20 (default) or  \nf: full (1-65535) or \ne.g. '23,53,80,443' or \ne.g. '100-200'")
}

func worker(target string, ports, results chan int) {
	for p := range ports {
		socket := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.Dial("tcp", socket)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

// scan ports and give back openports
func scanPorts(threads int, target string, ports []int) []int {
	portchan := make(chan int, scan.Threads)
	//portchan := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(target, portchan, results)
	}

	go func() {
		for _, p := range ports {
			portchan <- p
		}
	}()

	for i := 0; i < len(ports); i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(portchan)
	close(results)
	sort.Ints(openports)

	return openports
}

// parse ports from command-line-params and give back array of ports to scan
func parsePorts(portParam string) []int {
	var ports []int
	var portFields []string

	// portscan with top20
	if strings.Compare(portParam, "top20") == 0 {
		ports = top20
		return ports
	}

	// portscann all ports
	if strings.Compare(portParam, "full") == 0 {
		ports = getAllPorts()
		return ports
	}

	// check for multiple comma-separated ports
	if strings.ContainsRune(portParam, 44) {
		// split after each ','
		portFields = strings.Split(portParam, ",")
	} else {
		portFields = append(portFields, portParam)
	}

	for c, field := range portFields {
		// check if ports where specified as range with '-'
		if strings.ContainsRune(field, 45) {
			var portStart, portEnd int

			portSplit := strings.Split(field, "-")
			if len(portSplit) > 2 || portSplit[0] >= portSplit[1] {
				fmt.Println("[!] something wrong with specified ports")
				os.Exit(1)
			} else {
				for c, p := range portSplit {
					port, err := strconv.Atoi(p)
					if err != nil {
						fmt.Println("[!] can't convert port parameter: ", portFields[c])
						os.Exit(1)
					}
					ports = append(ports, port)
					if c == 0 {
						portStart = port
					} else {
						portEnd = port
					}
				}

				for i := portStart + 1; i < portEnd; i++ {
					ports = append(ports, i)
				}
			}

		} else {
			// Check if port param is not empty
			if len(field) > 0 {
				port, err := strconv.Atoi(field)
				if err != nil {
					fmt.Println("[!] can't convert port parameter: ", portFields[c])
					os.Exit(1)
				}

				ports = append(ports, port)
			}

		}
	}

	return ports
}

// parse targets from target param
func parseTargets(targetField string) ([]string, error) {
	var targets []string

	// Check for CIDR notation
	if strings.ContainsRune(targetField, 47) {
		target, network, err := net.ParseCIDR(targetField)
		if err != nil {
			return nil, err
		}

		// generate ip list from network address and network mask
		for target := target.Mask(network.Mask); network.Contains(target); inc(target) {
			targets = append(targets, target.String())
		}

		// remove network and broadcast address from list
		return targets[1 : len(targets)-1], nil

	} else {
		targets = append(targets, targetField)
		return targets, nil
	}
}

func inc(target net.IP) {
	for j := len(target) - 1; j >= 0; j-- {
		target[j]++
		if target[j] > 0 {
			break
		}
	}
}

// return full port list
func getAllPorts() []int {
	var ports []int

	for i := 1; i < 65536; i++ {
		ports = append(ports, i)
	}

	return ports
}

func main() {
	var targets []string

	// Parse command-line flags
	flag.Parse()

	// Parse targets from commandline
	targets = flag.Args()
	for _, targetField := range targets {
		t, err := parseTargets(targetField)
		if err != nil {
			fmt.Println("[!] error parsing cidr notation")
			os.Exit(1)
		}
		for _, target := range t {
			scan.Targets = append(scan.Targets, target)
		}
	}

	// Check if one or more targets are specified
	if len(scan.Targets) < 1 {
		fmt.Println("[!] No target specified\n")
		fmt.Println("Use ./portscanner <target>\nsingle target: ./portscanner 127.0.0.1\nmultiple targets: ./portscanner 192.168.0.1 192.168.0.2\nnetwork target: ./portscanner 10.10.10.0/24\n")
		os.Exit(0)
	}

	// parse port parameters
	scan.Ports = parsePorts(portParam)

	for c, target := range scan.Targets {
		fmt.Printf("\n[*] Scanning %d ports on target #%d (from %d) %s ...\n", len(scan.Ports), c+1, len(scan.Targets), target)
		// start timer for target
		start := time.Now()

		openports := scanPorts(scan.Threads, target, scan.Ports)

		fmt.Printf("\n---\nscan report for #%d (%s)\n\n", c+1, target)
		for _, port := range openports {
			fmt.Printf("%d/TCP open\n", port)
		}
		// get time
		duration := time.Since(start)

		fmt.Printf("---\n\n[i] finished target #%d in %s\n", c+1, duration)
	}

	// counter i starts with 0. +1 makes better human readability
	fmt.Printf("\n[i] finished %d targets.\n\n", len(scan.Targets))

}
