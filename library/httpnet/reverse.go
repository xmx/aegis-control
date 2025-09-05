package httpnet

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/xmx/aegis-common/contract/problem"
)

func NewReverse(trap http.RoundTripper) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetXForwarded()
		},
		Transport: trap,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			code := http.StatusBadGateway
			pb := &problem.Details{
				Host:     r.Host,
				Type:     r.Host,
				Status:   code,
				Detail:   err.Error(),
				Instance: r.URL.Path,
				Method:   r.Method,
				Datetime: time.Now().UTC(),
			}
			if ae, ok := err.(*net.AddrError); ok {
				pb.Detail = ae.Err
			}
			_ = pb.JSON(w)
		},
	}
}
