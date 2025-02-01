package main

import (
	"log"

	"github.com/haoran-mc/sniffer/packet"
	"github.com/haoran-mc/sniffer/setting"
)

func main() {
	if err := packet.Listen(setting.App.Nic); err != nil {
		log.Println(err.Error())
	}
}
