package victoria

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/xmx/metrics"
)

func NewProxy(cfg func(ctx context.Context) (pushURL string, opts *metrics.PushOptions, err error)) http.Handler {
	return &proxyMetrics{cfg: cfg}
}

type proxyMetrics struct {
	cfg func(ctx context.Context) (pushURL string, opts *metrics.PushOptions, err error)
}

func (pm *proxyMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pushURL, opts, err := pm.cfg(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pu, err := url.Parse(pushURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body io.Reader = r.Body
	if ce := r.Header.Get("Content-Encoding"); strings.EqualFold(ce, "gzip") {
		gzr, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer gzr.Close()
		body = gzr
	}
	raw, err := io.ReadAll(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if labels := opts.ExtraLabels; labels != "" {
		raw, _ = pm.addExtraLabels(nil, raw, labels)
	}

	var rbody io.Reader
	enableCompression := !opts.DisableCompression
	if enableCompression {
		buf := new(bytes.Buffer)
		gzw := gzip.NewWriter(buf)
		if _, err = gzw.Write(raw); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = gzw.Close(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rbody = buf
	} else {
		rbody = bytes.NewReader(raw)
	}

	method := opts.Method
	if method == "" {
		method = http.MethodGet
	}
	req, err := http.NewRequestWithContext(ctx, method, pushURL, rbody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, h := range opts.Headers {
		key, val, found := strings.Cut(h, ":")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		req.Header.Add(key, val)
	}
	if enableCompression {
		req.Header.Set("Content-Encoding", "gzip")
	}

	prx := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(pu)
		},
	}
	prx.ServeHTTP(w, req)
}

// addExtraLabels
func (pm *proxyMetrics) addExtraLabels(dst, src []byte, extraLabels string) ([]byte, bool) {
	var bashBytes = []byte("#")

	for len(src) > 0 {
		var line []byte
		n := bytes.IndexByte(src, '\n')
		if n >= 0 {
			line = src[:n]
			src = src[n+1:]
		} else {
			line = src
			src = nil
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			// Skip empy lines
			continue
		}
		if bytes.HasPrefix(line, bashBytes) {
			// Copy comments as is
			dst = append(dst, line...)
			dst = append(dst, '\n')
			continue
		}
		n = bytes.IndexByte(line, '{')
		if n >= 0 {
			dst = append(dst, line[:n+1]...)
			dst = append(dst, extraLabels...)
			dst = append(dst, ',')
			dst = append(dst, line[n+1:]...)
		} else {
			n = bytes.LastIndexByte(line, ' ')
			if n < 0 {
				return nil, false
			}
			dst = append(dst, line[:n]...)
			dst = append(dst, '{')
			dst = append(dst, extraLabels...)
			dst = append(dst, '}')
			dst = append(dst, line[n:]...)
		}
		dst = append(dst, '\n')
	}
	return dst, true
}
