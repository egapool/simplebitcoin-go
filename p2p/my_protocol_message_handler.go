package p2p

import "fmt"

type MyProtocolMessageHandler struct {
}

func (h *MyProtocolMessageHandler) handleMessage(msg string) {
	fmt.Println("MyProtocolMessageHandler received ", msg)
}
