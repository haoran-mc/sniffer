package main

import (
	"github.com/haoran-mc/sniffer/input"
	"github.com/haoran-mc/sniffer/mock"
	"github.com/haoran-mc/sniffer/reassembler"
	"github.com/haoran-mc/sniffer/setting"
)

func main() {
	go input.Listen(setting.App.Nic)
	go mock.MockServerStart()
	reassembler.StreamParserServerStart()
}
