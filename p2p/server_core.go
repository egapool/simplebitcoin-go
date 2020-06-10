package p2p

import (
	"log"
	"net"
	"os"
)

const (
	STATE_INIT                 int = 0
	STATE_STANDBY              int = 1
	STATE_CONNECTED_TO_NETWORK int = 2
	STATE_SHUTTING_DOWN        int = 3
)

type ServerCore struct {
	serverState int
	cm          ConnectionManager
	myPort      string
	myIp        string
	coreNode    Host
	mpmh        MyProtocolMessageHandler
}

// NewGenesiCore は始原のcoreノードを生成する
func NewGenesiCore(port string) ServerCore {
	log.Println("Initializing Genesis server...")
	core := ServerCore{
		myPort: port,
	}
	core.myIp = core.getMyIp()
	log.Printf("Server IP address is set to ... %s\n", core.myIp)
	core.cm = NewConnectionManager(core.myIp, core.myPort, core.handleMessage)
	return core
}

func NewServerCore(port string, coreNodeHost string, coreNodePort string) ServerCore {
	log.Println("Initializing server...")
	core := ServerCore{
		myPort: port,
	}
	core.myIp = core.getMyIp()
	log.Printf("Server IP address is set to ... %s\n", core.myIp)
	core.cm = NewConnectionManager(core.myIp, core.myPort, core.handleMessage)
	core.coreNode = Host{Hostname: coreNodeHost, Port: coreNodePort}
	return core
}

func (c *ServerCore) handleMessage(result string, reason int, cmd int, peer_port string, payload string, host Host) {
	if cmd == MSG_NEW_TRANSACTION {
	} else if cmd == MSG_NEW_BLOCK {

	} else if cmd == RSP_FULL_CHAIN {

	} else if cmd == MSG_ENHANCED {
		c.mpmh.handleMessage(payload)
	}

}

// Start is start core server
func (c *ServerCore) Start() {
	c.serverState = STATE_STANDBY
	c.cm.Start()
}

func (c *ServerCore) JoinNetwork() {
	c.serverState = STATE_CONNECTED_TO_NETWORK
	c.cm.JoinNetwork(c.coreNode)
}

func (c *ServerCore) Shutdown() {
	c.serverState = STATE_SHUTTING_DOWN
	log.Println("Shutdown server...")
	c.cm.ConnectionClose()
}

func (c *ServerCore) getMyCurrentState() int {
	return c.serverState
}

func (c *ServerCore) getMyIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
