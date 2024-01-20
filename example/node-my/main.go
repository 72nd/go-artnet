package main

import (
	"fmt"
	"net"
	"time"

	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

func Description(desc string) artnet.NodeOption {
	return func(n *artnet.Node) error {
		n.Config.Description = n.Config.Description
		return nil
	}
}

func Input(data []artnet.InputPort) artnet.NodeOption {
	return func(n *artnet.Node) error {
		n.Config.InputPorts = data
		return nil
	}
}

func main() {
	artsubnet := "2.0.0.0/8"
	_, cidrnet, _ := net.ParseCIDR(artsubnet)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("error getting ips: %s\n", err)
	}

	var ip net.IP
	fmt.Printf("%+v", addrs)

	for _, addr := range addrs {
		ip = addr.(*net.IPNet).IP
		if cidrnet.Contains(ip) {
			break
		}
	}

	ip = net.ParseIP("192.168.3.79")

	log := artnet.NewDefaultLogger()
	n := artnet.NewNode("Shelly Gateway", code.StNode, ip, log)
	n.Config.Description = "YOLO"
	customBroadcastAddr := net.UDPAddr{
		IP:   net.IPv4(192, 168, 3, 255),
		Port: 6454,
	}
	n.SetOption(artnet.NodeBroadcastAddress(customBroadcastAddr))
	n.SetOption(Input([]artnet.InputPort{
		{
			Address: artnet.Address{
				Net:    0,
				SubUni: 0,
			},
			Type:   0xc0,
			Status: 0,
		},
	}))
	n.Config.OutputPorts = []artnet.OutputPort{
		{
			Address: artnet.Address{
				Net:    0,
				SubUni: 0,
			},
			Type:   0xc0,
			Status: 0,
		},
	}
	n.RegisterCallback(code.OpOutput, func(p packet.ArtNetPacket) {
		fmt.Println("JAAA")
		pkg, ok := p.(*packet.ArtCommandPacket)
		if !ok {
			return
		}
		fmt.Println(pkg.Header)
	})
	fmt.Printf("%+v", n.Config)
	n.Start()

	for {
		time.Sleep(time.Second)
	}
}
