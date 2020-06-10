package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/aerialpartners/study-cryptocurrency/libs"
)

type ConnectionManager4Edge struct {
	host         string
	port         string
	myCoreHost   string
	myCorePort   string
	coreNodeList CoreNodeList
}

// NewConnecitonManager4Edge is managing connection of edge
func NewConnecitonManager4Edge(myIp string, myPort string, myCoreIp string, myCorePort string) (cm ConnectionManager4Edge) {
	cm = ConnectionManager4Edge{
		host:         myIp,
		port:         myPort,
		myCoreHost:   myCoreIp,
		myCorePort:   myCorePort,
		coreNodeList: NewCoreNodeList(),
	}
	return
}

func (cm *ConnectionManager4Edge) start() {
	// code
	cm.waitForAccess()
}

func (cm *ConnectionManager4Edge) connectToCoreNode() {
	cm.connectToP2PNW(cm.myCoreHost, cm.myCorePort)
}

func (cm *ConnectionManager4Edge) sendMsg(host Host, msg Message) {
	fmt.Println("Sending...", msg)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host.Service())
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println("Connection failed for host: ", host)
		cm.coreNodeList.remove(host)
		// todo 接続失ったら次のホストへ行く実装
	} else {
		msg_json, _ := json.Marshal(msg)
		_, _ = conn.Write(msg_json)

	}
}

func (cm *ConnectionManager4Edge) connectToP2PNW(host string, port string) {
	msg := NewMessage(MSG_ADD_AS_EDGE, cm.port, "")
	cm.sendMsg(Host{host, port}, msg)
}

func (cm *ConnectionManager4Edge) waitForAccess() {
	service := ":" + cm.port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	libs.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	libs.CheckError(err)
	fmt.Println("waiting for the connection ..." + cm.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		} else {
			fmt.Println("connected by .. ", conn.RemoteAddr().String())
		}
		go func() {
			// set 2 minutes timeout
			conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
			// TODO バイナリがすべて完結するまでfor回す必要あり
			request := make([]byte, 1024)
			defer conn.Close() // close connection before exit
			readLen, err := conn.Read(request)
			if err != nil {
				fmt.Println("Request read error ", err, readLen)
				return
			}
			cm.handleMessage(conn, request[:readLen])
		}()
	}
}
func (cm *ConnectionManager4Edge) handleMessage(conn net.Conn, msg []byte) {
	result, reason, cmd, peer_port, payload := parse(msg)
	fmt.Println(result, reason, cmd, peer_port, payload)
	if result == "error" {
		if reason == ERR_PROTOCOL_UNMATCH {
			fmt.Println("Error: Protocol name is not matched")
			return
		} else if reason == ERR_VERSION_UNMATCH {
			fmt.Println("Error: Protocol version is not matched")
			return
		}
	} else if result == "ok" {
		if reason == OK_WITHOUT_PAYLOAD {
			if cmd == MSG_PING {
				// pass
			} else {
				fmt.Println("Edge node does not have functions for this message!")
			}
		} else if reason == OK_WITH_PAYLOAD {
			if cmd == MSG_CORE_LIST {
				fmt.Println("Refresh the core node list...")
				var nodeList [][2]string
				err := json.Unmarshal([]byte(payload), &nodeList)
				if err != nil {
					fmt.Println("JSON parse error ! msg_core_list ", err)
				}
				for _, host := range nodeList {
					cm.coreNodeList.add(Host{host[0], host[1]})
				}
			} else {
				// callback
			}
		}
	} else {
		fmt.Println("Unexpected status", result, reason)
	}
}
func (cm *ConnectionManager4Edge) sendPing() {
	host := Host{cm.myCoreHost, cm.myCorePort}
	message := NewMessage(MSG_PING, cm.port, "")
	cm.sendMsg(host, message)
}

func (cm *ConnectionManager4Edge) getMessageText(msgType int, payload string) Message {
	return NewMessage(msgType, cm.port, payload)
}
