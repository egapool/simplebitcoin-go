package main

import "github.com/egapool/simplebitcoin-go/p2p"

func main() {
	server := p2p.NewGenesiCore("50082")
	server.Start()
}
