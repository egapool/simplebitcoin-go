package p2p

type Host struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
}

// Service はIPとPortの連結文字列を返す
func (host *Host) Service() string {
	return host.Hostname + ":" + host.Port
}
