package model

import "time"

type NodeNetwork struct {
	Name  string   `json:"name"          bson:"name"`
	Index int      `json:"index"         bson:"index"`
	MTU   int      `json:"mtu"           bson:"mtu"`
	IPv4  []string `json:"ipv4,omitzero" bson:"ipv4,omitempty"`
	IPv6  []string `json:"ipv6,omitzero" bson:"ipv6,omitempty"`
	MAC   string   `json:"mac,omitzero"  bson:"mac,omitempty"`
}

type NodeNetworks []*NodeNetwork

type TunnelStat struct {
	ConnectedAt    time.Time `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`
	DisconnectedAt time.Time `json:"disconnected_at,omitzero" bson:"disconnected_at,omitempty"`
	KeepaliveAt    time.Time `json:"keepalive_at,omitzero"    bson:"keepalive_at,omitempty"`
	Protocol       string    `json:"protocol,omitzero"        bson:"protocol,omitempty"`
	Subprotocol    string    `json:"subprotocol,omitzero"     bson:"subprotocol,omitempty"`
	LocalAddr      string    `json:"local_addr,omitzero"      bson:"local_addr,omitempty"`
	RemoteAddr     string    `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`
	ReceiveBytes   uint64    `json:"receive_bytes,omitzero"   bson:"receive_bytes,omitempty"`  // broker/agent 为主体
	TransmitBytes  uint64    `json:"transmit_bytes,omitzero"  bson:"transmit_bytes,omitempty"` // broker/agent 为主体
}

type ExecuteStat struct {
	Goos       string   `json:"goos,omitzero"       bson:"goos,omitempty"`
	Goarch     string   `json:"goarch,omitzero"     bson:"goarch,omitempty"`
	PID        int      `json:"pid,omitzero"        bson:"pid,omitempty"`
	Args       []string `json:"args,omitzero"       bson:"args,omitempty"`
	Hostname   string   `json:"hostname,omitzero"   bson:"hostname,omitempty"`
	Workdir    string   `json:"workdir,omitzero"    bson:"workdir,omitempty"`
	Executable string   `json:"executable,omitzero" bson:"executable,omitempty"`
	Username   string   `json:"username,omitzero"   bson:"username,omitempty"`
	UID        string   `json:"uid,omitzero"        bson:"uid,omitempty"`
}
