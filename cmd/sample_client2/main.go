package main

import (
	"encoding/json"
	"time"

	"github.com/egapool/simplebitcoin-go/p2p"
)

func main() {
	client := p2p.NewClientCore("50095", "127.0.0.1", "50082")
	go func() {
		time.Sleep(time.Second * 5)
		msg := map[string]string{"message": "test"}
		json, _ := json.Marshal(msg)
		client.SendMessageToMyCoreNode(p2p.MSG_ENHANCED, string(json))
	}()
	client.Start()
}
