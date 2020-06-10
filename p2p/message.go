package p2p

import (
	"encoding/json"
	"fmt"
)

const (
	ProtocolName           string = "simple_bitcoin_protocol"
	MY_VERSION             string = "0.1.0"
	MSG_ADD                int    = 0
	MSG_REMOVE             int    = 1
	MSG_CORE_LIST          int    = 2
	MSG_REQUEST_CORE_LIST  int    = 3
	MSG_PING               int    = 4
	MSG_ADD_AS_EDGE        int    = 5
	MSG_REMOVE_EDGE        int    = 6
	MSG_NEW_TRANSACTION    int    = 7
	MSG_NEW_BLOCK          int    = 8
	MSG_REQUEST_FULL_CHAIN int    = 9
	RSP_FULL_CHAIN         int    = 10
	MSG_ENHANCED           int    = 11

	ERR_PROTOCOL_UNMATCH int = 0
	ERR_VERSION_UNMATCH  int = 1
	OK_WITH_PAYLOAD      int = 2
	OK_WITHOUT_PAYLOAD   int = 3
	JSON_UNMARSHAL_ERROR int = 4
)

// Message is
type Message struct {
	Protocol string `json:"protocl"`
	Version  string `json:"version"`
	MsgType  int    `json:"msg_type"`
	MyPort   string `json:"my_port"`
	Payload  string `json:"payload"`
}

func NewMessage(msgType int, port string, payload string) Message {
	return Message{
		Protocol: ProtocolName,
		Version:  MY_VERSION,
		MsgType:  msgType,
		MyPort:   port,
		Payload:  payload,
	}
}

func parse(msg []byte) (string, int, int, string, string) {
	fmt.Println(string(msg))
	data := new(Message)

	if err := json.Unmarshal(msg, data); err != nil {
		fmt.Println("Failed Unmarshal", err)
		return "error", JSON_UNMARSHAL_ERROR, 500, "0", ""
	}
	if data.Protocol != ProtocolName {
		return "error", ERR_PROTOCOL_UNMATCH, 500, "0", ""
	} else if data.Version > MY_VERSION {
		return "error", ERR_VERSION_UNMATCH, 500, "0", ""
	} else if data.MsgType == MSG_CORE_LIST || data.MsgType == MSG_NEW_TRANSACTION || data.MsgType == MSG_NEW_BLOCK || data.MsgType == RSP_FULL_CHAIN || data.MsgType == MSG_ENHANCED {
		return "ok", OK_WITH_PAYLOAD, data.MsgType, data.MyPort, data.Payload
	} else {
		return "ok", OK_WITHOUT_PAYLOAD, data.MsgType, data.MyPort, ""
	}
}
