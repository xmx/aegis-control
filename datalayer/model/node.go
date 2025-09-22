package model

type NodeNetwork struct {
	Name  string   `json:"name"          bson:"name"`
	Index int      `json:"index"         bson:"index"`
	MTU   int      `json:"mtu"           bson:"mtu"`
	IPv4  []string `json:"ipv4,omitzero" bson:"ipv4,omitempty"`
	IPv6  []string `json:"ipv6,omitzero" bson:"ipv6,omitempty"`
	MAC   string   `json:"mac,omitzero"  bson:"mac,omitempty"`
}

type NodeNetworks []*NodeNetwork
