package p2p

// EdgeNodeList はEdgeノードを管理するための構造体
type EdgeNodeList struct {
	list map[Host]bool
}

//type Edge struct {
//	Host
//}

func NewEdgeNodeList() EdgeNodeList {

	return EdgeNodeList{list: make(map[Host]bool)}
}

func (l *EdgeNodeList) add(edge Host) {
	l.list[edge] = true
}

func (l *EdgeNodeList) remove(edge Host) {
	_, ok := l.list[edge]
	if ok {
		delete(l.list, edge)
	}
}

func (l *EdgeNodeList) overwrite() {}
