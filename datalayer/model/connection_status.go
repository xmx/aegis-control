package model

import "time"

type ConnectionStatus struct {
	ConnectedAt    time.Time `json:"connected_at,omitzero"    bson:"connected_at,omitempty"`
	DisconnectedAt time.Time `json:"disconnected_at,omitzero" bson:"disconnected_at,omitempty"`
	KeepaliveAt    time.Time `json:"keepalive_at,omitzero"    bson:"keepalive_at,omitempty"`
	Protocol       string    `json:"protocol,omitzero"        bson:"protocol,omitempty"`
	Subprotocol    string    `json:"subprotocol,omitzero"     bson:"subprotocol,omitempty"`
	RemoteAddr     string    `json:"remote_addr,omitzero"     bson:"remote_addr,omitempty"`
}
