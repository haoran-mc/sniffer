package input

import (
	"io"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/haoran-mc/sniffer/setting"
)

func Listen(nic string) {
	// 获取网卡信息
	iface, err := net.InterfaceByName(nic)
	if err != nil {
		log.Fatalf("NIC %s not found, error: %v", setting.App.Nic, err)
	}
	log.Printf("NIC: %s, MTU: %d", nic, iface.MTU)

	// 打开设备监听
	handle, err := pcap.OpenLive(
		nic,                // Nic
		9000,               // SnapLen
		true,               // Promisc
		1*time.Millisecond, // Timeout
	)
	if err != nil {
		log.Fatalf("OpenLive %s error: %v", setting.App.Nic, err)
	}

	// 设置过滤
	if err := handle.SetBPFFilter(setting.App.Bpf); err != nil {
		log.Fatalf("Set BPF Filter error: %v", err)
	}

	if handle.LinkType() != layers.LinkTypeEthernet {
		log.Fatal("Handle link type is not ethernet")
	}

	go func() {
		defer handle.Close()
		pcktChanLen := cap(pktChan)

		for {
			data, _, err := handle.ReadPacketData()
			if err == nil {
				if len(pktChan) < pcktChanLen {
					pktChan <- &packet{data}
				} // else queue drop
				continue
			}
			if enext, ok := err.(pcap.NextError); ok && enext == pcap.NextErrorTimeoutExpired {
				continue
			}
			if eno, ok := err.(syscall.Errno); ok && eno.Temporary() {
				continue
			}
			if enet, ok := err.(*net.OpError); ok && (enet.Temporary() || enet.Timeout()) {
				continue
			}
			if err == io.EOF || err == io.ErrClosedPipe {
				log.Fatalf("stopped reading from %s interface with error: %s\n", nic, err)
			}

			log.Fatalf("stopped reading from %s interface with error: %s\n", nic, err)
		}
	}()
}
