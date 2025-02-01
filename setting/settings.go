package setting

import (
	"flag"
)

type AppSetting struct {
	Nic string
	Bpf string
}

var App = new(AppSetting)

func init() {
	flag.StringVar(&App.Nic, "i", "en0", "NIC name")
	flag.StringVar(&App.Bpf, "bpf", "", "BPF(Berkeley Packet Filter) String")
	flag.Parse()
}
