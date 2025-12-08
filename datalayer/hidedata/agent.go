package hidedata

type Agent struct {
	Protocols []string `json:"protocols"` // 连接协议 udp tcp，一般留空即可。
	Addresses []string `json:"addresses"` // 连接的 broker 地址
	Offset    int64    `json:"offset"`
}
