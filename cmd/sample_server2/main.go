package main

import "github.com/egapool/simplebitcoin-go/p2p"

func main() {
	server := p2p.NewServerCore("50083", "192.168.1.3", "50082")
	server.JoinNetwork()
	server.Start()
}
