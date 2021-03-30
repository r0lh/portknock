<!-- PROJECT HEADER -->
[![License](https://img.shields.io/badge/License-BSD%203--Clause-orange.svg)](https://github.com/r0lh/portknock/LICENSE.txt)

# portknock
	
Portknock is a simple portscanner written in Go.

<!-- TABLE OF CONTENTS -->
<details open="open">
	<summary><h2 style="display: inline-block">Table of Contents</h2></summary>
	<ol>
	  <li>
	    <a href="#about-the-project">About The Project</a>
	    <ul>
	      <li><a href="built-with">Built With</a></li>
	    </ul>
	  </li>
	  <li>
	    <a href="#getting-started">Getting Started</a>
	    <ul>
	      <li><a href="#prerequisites">Prerequisites</a></li>
	      <li><a href="#installation">Installation</a></li>
	    </ul>
	  </li>
	  <li><a href="#usage">Usage</a></li>
	  <li><a href="#license">License</a></li>
	  <li><a href="#contact">Contact</a></li>
	  <li><a href="#bugs-todo">Bugs / ToDo</a></li>
	</ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

Portknock is a simple portscanner built with Go. Go give the advantage of a static linked binary, so the portscanner can be built for different platforms and no other dependencies for the execution on the target are needed.

### Built With
* [The Go Programming Language](https://golang.org)

<!-- GETTING STARTED -->
## Getting Started
To use portknock on your system follow these steps.

### Prerequisites
Install Go on your system. Follow the official [Go Docs](https://golang.org/doc/install) for your platform.

### Installation
1. Clone the repo
```sh
git clone https://github.com/r0lh/portknock.git
```
2. Build the binary
```sh
cd portknock && go build portknock.go
```
3. Make the portknock binary availabe for all users on your system (Unix/BSD/Linux/Mac OS X)
```sh
sudo cp ./portknock /usr/local/bin
```
<!-- USAGE EXAMPLES -->
## Usage

Make a portscan on a single target. Portknock scans per default the top20 most scanned ports from nmap database.
```sh
portknock 127.0.0.1
```

Make a portscan on multiple targets.
```sh
portknock 192.168.10.1 192.168.168.10.105
```

Make a portscan on a network with [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#CIDR_notation)
```sh
portknock 192.168.0.0/24
```

Make a fullscan (port 1-65535)
```sh
portknock -p full  10.10.10.1
```

Scan specified ports
```sh
portknock -p 80,443 192.168.0.1
```

Set more or less concurrent threads
```sh
portknow -t 30 192.168.0.10 
```

<!-- LICENSE -->
## License

Distributed under the BSD 3-clause license. See `LICENSE` for more information.

<!-- CONTACT -->
## Contact

Roland Heimpoldinger - roland@heimpoldinger.net

Project Link: [https://github.com/r0lh/portknock](https://github.com/r0lh/portknock)

<!-- BUGS / TODO -->
## Bugs / ToDo

* portknock can't handle URLs yet
* portknock don't ping targets before scanning, so even not reachable targets are scanned
