package main

import "github.com/aerialpartners/study-cryptocurrency/p2p"

func main() {
	server := p2p.NewGenesiCore("50082")
	server.Start()
}
