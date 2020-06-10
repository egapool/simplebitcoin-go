package p2p

import (
	"encoding/json"
	"log"
)

// CoreNodeList はコアノードのリスト
type CoreNodeList struct {
	list map[Host]bool
}

func NewCoreNodeList() CoreNodeList {
	return CoreNodeList{list: make(map[Host]bool)}
}

func (l *CoreNodeList) add(host Host) {
	l.list[host] = true
	log.Println("Current Core List: ", l.toJson())
}

func (l *CoreNodeList) remove(host Host) {
	_, ok := l.list[host]

	if ok {
		delete(l.list, host)
	}
}

func (l *CoreNodeList) toJson() string {
	var tmp [][2]string
	for host, _ := range l.list {
		p := [2]string{host.Hostname, host.Port}
		tmp = append(tmp, p)
	}
	j, _ := json.Marshal(tmp)
	return string(j)
}
