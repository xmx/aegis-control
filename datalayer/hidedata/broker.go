package hidedata

type Broker struct {
	Protocols []string `json:"protocols"`
	Addresses []string `json:"addresses"`
	Secret    string   `json:"secret"`
	Offset    int64    `json:"offset"`
}
