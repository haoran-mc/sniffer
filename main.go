package main

import (
	"github.com/haoran-mc/sniffer/input"
	"github.com/haoran-mc/sniffer/setting"
)

func main() {
	go input.Listen(setting.App.Nic)

	select {}
}
