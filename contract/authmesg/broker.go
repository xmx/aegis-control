package authmesg

import "net/http"

type BrokerToServerRequest struct {
	Secret     string   `json:"secret"              validate:"required,gte=10,lte=100"`
	Goos       string   `json:"goos"                validate:"required,oneof=darwin dragonfly illumos ios js wasip1 linux android solaris freebsd nacl netbsd openbsd plan9 windows aix"`
	Goarch     string   `json:"goarch"              validate:"required,oneof=386 amd64 arm arm64 loong64 mips mipsle mips64 mips64le ppc64 ppc64le riscv64 s390x sparc64 wasm"`
	PID        int      `json:"pid,omitzero"`
	Args       []string `json:"args,omitzero"`
	Hostname   string   `json:"hostname,omitzero"`
	Workdir    string   `json:"workdir,omitzero"`
	Executable string   `json:"executable,omitzero"`
	Username   string   `json:"username,omitzero"`
	UID        string   `json:"uid,omitzero"`
}

type ServerToBrokerResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message,omitzero"`
	Config  BrokerInitialConfig `json:"config"`
}

type BrokerInitialConfig struct {
	URI string `json:"uri"` // 数据库连接地址
}

func (r ServerToBrokerResponse) Succeed() bool {
	return r.Code >= http.StatusOK && r.Code < http.StatusMultipleChoices
}
