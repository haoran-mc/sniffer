package main

import (
	"github.com/haoran-mc/sniffer/input"
	"github.com/haoran-mc/sniffer/output/db"
	"github.com/haoran-mc/sniffer/proxy"
	"github.com/haoran-mc/sniffer/replay"
	"github.com/haoran-mc/sniffer/setting"
)

func main() {
	go input.Listen(setting.App.Nic)
	go replay.StartResponseServer()

	db.InitClickhouse()
	proxy.StartServer()
}
