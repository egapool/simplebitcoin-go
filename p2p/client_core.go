package p2p

import (
	"log"
	"net"
	"os"
)

const (
	STATE_ACTIVE = 4
)

type ClientCore struct {
	clientState int
	myIp        string
	myPort      string
	myCoreIp    string
	myCorePort  string
	cm          ConnectionManager4Edge
}

func NewClientCore(myPort string, coreHost string, corePort string) ClientCore {
	log.Println("Initializing server...")
	core := ClientCore{
		myPort:      myPort,
		myCoreIp:    coreHost,
		myCorePort:  corePort,
		clientState: STATE_INIT,
	}
	core.myIp = core.getMyIp()
	core.cm = NewConnecitonManager4Edge(core.myIp, core.myPort, coreHost, corePort)
	return core
}

func (c *ClientCore) Start() {
	c.clientState = STATE_ACTIVE
	c.cm.connectToCoreNode()
	c.cm.start()
}

func (c *ClientCore) shutdown() {
	c.clientState = STATE_SHUTTING_DOWN
	log.Println("Shutdown edge node...")
	// connection close
}

func (c *ClientCore) SendMessageToMyCoreNode(msgType int, msg string) {
	msgTxt := c.cm.getMessageText(msgType, msg)
	c.cm.sendMsg(Host{c.myCoreIp, c.myCorePort}, msgTxt)

}
func (c *ClientCore) getMyIp() string {
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
