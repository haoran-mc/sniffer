package packet

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/haoran-mc/sniffer/setting"
)

func Listen(nic string) error {
	// 获取网卡信息
	iface, err := net.InterfaceByName(nic)
	if err != nil {
		return fmt.Errorf("nic %s not found, err: %v", setting.App.Nic, err)
	}
	log.Printf("NIC: %s, MTU: %d", nic, iface.MTU)

	// 打开设备监听
	handle, err := pcap.OpenLive(
		setting.App.Nic, // 网卡名
		1024*1024,       // snaplen
		true,            // 混杂模式
		pcap.BlockForever,
	)
	if err != nil {
		return fmt.Errorf("openLive %s err: %v", setting.App.Nic, err)
	}
	defer handle.Close()

	// 设置过滤器
	if err := handle.SetBPFFilter(setting.App.Bpf); err != nil {
		return fmt.Errorf("set bpf filter: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	packetSource.Lazy = true

	for packet := range packetSource.Packets() {
		go analysePacket(packet)
	}
	return nil
}

func analysePacket(packet gopacket.Packet) {
	// ipv4
	if ip4Layer := packet.Layer(layers.LayerTypeIPv4); ip4Layer != nil {
		ip4 := ip4Layer.(*layers.IPv4)

		// tcp
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)

			// is HTTP request or response
			if len(tcp.Payload) > 4 {
				if bytes.Equal(tcp.Payload[:4], []byte("HTTP")) {
					fmt.Printf("%s	Respond			From %s:%d to %s:%d\n", tcp.Payload[:4], ip4.SrcIP, tcp.SrcPort, ip4.DstIP, tcp.DstPort)
				} else if bytes.Equal(tcp.Payload[:4], []byte("GET ")) || bytes.Equal(tcp.Payload[:4], []byte("POST")) {

					// find the path of the request
					i1 := -1
					i2 := -1
					for index, value := range tcp.Payload {
						if i1 == -1 {
							if value == ' ' {
								i1 = index + 1
							}
						} else if i2 == -1 {
							if value == ' ' {
								i2 = index
								break
							}
						}
					}
					fmt.Printf("HTTP 	%s %s			From %s:%d to %s:%d\n", tcp.Payload[:4], tcp.Payload[i1:i2], ip4.SrcIP, tcp.SrcPort, ip4.DstIP, tcp.DstPort)
				}
			}
		}
	}
}
