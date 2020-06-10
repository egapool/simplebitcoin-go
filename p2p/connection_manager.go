package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/egapool/simplebitcoin-go/libs"
)

// ConnectionManager は Coreどうしの接続を管理する
type ConnectionManager struct {
	host         string
	port         string
	coreNodeList CoreNodeList
	edgeNodeList EdgeNodeList
	callback     func(result string, reason int, cmd int, peer_port string, payload string, host Host)
}

// NewConnectionManager is create new ConnectionManager
func NewConnectionManager(myIp string, myPort string, callback func(result string, reason int, cmd int, peer_port string, payload string, host Host)) ConnectionManager {
	fmt.Println("Initializing ConnectionManager...")
	cm := ConnectionManager{
		host:         myIp,
		port:         myPort,
		coreNodeList: NewCoreNodeList(),
		edgeNodeList: NewEdgeNodeList(),
		callback:     callback,
	}
	cm.addPeer(Host{myIp, myPort})
	return cm
}

// Start は待受を開始する際に呼び出される
func (cm *ConnectionManager) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cm.checkPeersConnection(ctx)

	cm.waitForAccess()
}

// JoinNetwork はユーザーが指定した既知のCoreノードへの接続
func (cm *ConnectionManager) JoinNetwork(host Host) {
	log.Println("Join Network to " + host.Hostname + ":" + host.Port)
	cm.connectToP2PNW(host)
}

func (cm *ConnectionManager) connectToP2PNW(host Host) {
	msg := NewMessage(MSG_ADD, cm.port, "")
	cm.SendMsg(host, msg)
}

// SendMsg は指定されたノードに対してメッセージを送信する
func (cm *ConnectionManager) SendMsg(host Host, msg Message) {
	log.Println("Send Message...to", host)
	log.Println(msg)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host.Service())
	log.Println("send message to ", tcpAddr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		cm.removePeer(host)
	} else {
		msg_json, _ := json.Marshal(msg)
		_, _ = conn.Write(msg_json)
	}
}

// SendMsgToAllPeer はCoreノードリストに登録されているすてべのノードに対して
// 同じメッセージをブロードキャストする
func (cm *ConnectionManager) SendMsgToAllPeer(msg Message) {
	fmt.Println("send_msg_to_all_peer was called!")
	fmt.Println(cm.coreNodeList.list)
	for host, _ := range cm.coreNodeList.list {
		if !cm.isSelf(host) {
			fmt.Println("message will be sent to ... ", host)
			cm.SendMsg(host, msg)
		}
	}
}

// ConnectionClose は終了前の処理としてソケットを閉じる
func (cm *ConnectionManager) ConnectionClose() {
	// TODO ちょっとよくわからない
}

func (cm *ConnectionManager) handleMessage(conn net.Conn, msg []byte) {
	result, reason, cmd, peer_port, payload := parse(msg)
	addr := conn.RemoteAddr().String()
	addr = strings.Split(addr, ":")[0]
	host := Host{Hostname: addr, Port: peer_port}
	fmt.Println(host)
	fmt.Println(result, reason, cmd, peer_port, payload)

	if result == "error" {
		if reason == ERR_PROTOCOL_UNMATCH {
			log.Println("Error: Protocl name is not matched")
			return
		} else if result == "error" && reason == ERR_VERSION_UNMATCH {
			log.Println("Error: Protocl version is not matched")
			return
		}
	}
	if result == "ok" {
		if reason == OK_WITHOUT_PAYLOAD {
			if cmd == MSG_ADD {
				log.Println("ADD node request was received!")
				cm.addPeer(host)
				if cm.isSelf(host) {
					return
				}
				message := NewMessage(MSG_CORE_LIST, cm.port, cm.coreNodeList.toJson())
				cm.SendMsgToAllPeer(message)
			} else if cmd == MSG_REMOVE {
				log.Printf("REMOVE request was received! from %s \n", addr)
				cm.removePeer(host)
			} else if cmd == MSG_REQUEST_CORE_LIST {
				log.Println("List for Core nodes was requested!")
				msg := NewMessage(MSG_CORE_LIST, cm.port, cm.coreNodeList.toJson())
				cm.SendMsg(host, msg)
			} else if cmd == MSG_ADD_AS_EDGE {
				// code
				cm.addEdgeNode(host)
				msg := NewMessage(MSG_CORE_LIST, cm.port, cm.coreNodeList.toJson())
				cm.SendMsg(host, msg)
			} else if cmd == MSG_REMOVE_EDGE {
				// code
				cm.removeEdgeNode(host)
			} else if cmd == MSG_PING {
			} else {
				cm.callback(result, reason, cmd, peer_port, payload, host)
				return
			}
		} else if reason == OK_WITH_PAYLOAD {
			if cmd == MSG_CORE_LIST {
				log.Println("refresh the core node list...")
				var nodeList [][2]string
				err := json.Unmarshal([]byte(payload), &nodeList)
				if err != nil {
					fmt.Println("JSON parse error ! msg_core_list ", err)
				}
				for _, host := range nodeList {
					cm.coreNodeList.add(Host{host[0], host[1]})
				}
			} else {
				cm.callback(result, reason, cmd, peer_port, payload, host)
				return
			}
		}
	}
}

func (cm *ConnectionManager) addPeer(host Host) {
	log.Println("Adding host: ", host)
	cm.coreNodeList.add(host)
}

func (cm *ConnectionManager) removePeer(host Host) {
	cm.coreNodeList.remove(host)
}

func (cm *ConnectionManager) checkPeersConnection(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("check_peers_connection was called")
			changed := false
			for host := range cm.coreNodeList.list {
				if !cm.isSelf(host) && !cm.isAlive(host) {
					changed = true
					delete(cm.coreNodeList.list, host)
					fmt.Println("Nonactive Node is deleted. Current Node is...", cm.coreNodeList.toJson())
				}
			}
			if changed {
				message := NewMessage(MSG_CORE_LIST, cm.port, cm.coreNodeList.toJson())
				fmt.Println(message)
				cm.SendMsgToAllPeer(message)
			}
		}
	}
}

// ConnectionManagerに生やす必要があるのか？
func (cm *ConnectionManager) waitForAccess() {
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
			//conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
			// TODO バイナリがすべて完結するまでfor回す必要あり
			request := make([]byte, 1024)
			defer conn.Close() // close connection before exit
			readLen, err := conn.Read(request)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Request read error ", err, readLen)
					return
				}
			}
			cm.handleMessage(conn, request[:readLen])
		}()
	}
}

func (cm *ConnectionManager) isAlive(host Host) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host.Service())
	if err != nil {
		return false
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		return false
	}
	msg := NewMessage(MSG_PING, cm.port, "")
	cm.SendMsg(host, msg)
	return true
}

func (cm *ConnectionManager) isSelf(host Host) bool {
	return cm.host == host.Hostname && cm.port == host.Port
}

func (cm *ConnectionManager) addEdgeNode(host Host) {
	cm.edgeNodeList.add(host)
}

func (cm *ConnectionManager) removeEdgeNode(host Host) {
	cm.edgeNodeList.remove(host)
}
