package main

import (
	"fmt"
	"log"
	"net"
)

type MyNet struct {
	IP           *net.IPNet
	Subnet       net.IP
	AvailableIPs []net.IP
}

func main() {
	mn := new(MyNet)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				//	ip = v.IP
				continue
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			fmt.Println("IP: ", ip)
			fmt.Println("Subnet: ", ip.Mask(ip.DefaultMask()))
			mn.IP = addr.(*net.IPNet)
			mn.Subnet = ip.Mask(ip.DefaultMask())
		}
	}
	fmt.Println(mn)
	mn.GetNetwork()
}

func (m *MyNet) GetNetwork() {
	for ip := m.Subnet; m.IP.Contains(ip); inc(ip) {
		fmt.Println(ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
