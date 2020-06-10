package main

import "github.com/aerialpartners/study-cryptocurrency/p2p"

func main() {
	client := p2p.NewClientCore("50095", "127.0.0.1", "50082")
	client.Start()
}
