package tlscert

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log/slog"
	"math/big"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Certifier interface {
	Match(*tls.ClientHelloInfo) (*tls.Certificate, error)
	Reset()
}

func NewCertPool(slow func(context.Context) ([]*tls.Certificate, error), log *slog.Logger) Certifier {
	return &certPool{
		slow: slow,
		log:  log,
	}
}

type certPool struct {
	slow  func(context.Context) ([]*tls.Certificate, error) // 惰性加载函数。
	log   *slog.Logger                                      // 日志输出
	mutex sync.Mutex
	cert  atomic.Pointer[certMap]
	self  atomic.Pointer[tls.Certificate] // 自签证书，兜底用。
}

func (cp *certPool) Match(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	sni := chi.ServerName
	attrs := []any{"sni", sni}
	cp.log.Debug("开始匹配合适的证书", attrs...)

	crtm := cp.cert.Load()
	if crtm == nil {
		cp.log.Debug("懒加载证书池", attrs...)
		crtm = cp.slowLoad()
	}
	if crt, err := crtm.Match(sni); crt != nil {
		cp.log.Debug("证书池中匹配到了合适的证书", attrs...)
		return crt, nil
	} else {
		args := append(attrs, "error", err)
		cp.log.Debug("证书池中未匹配到合适的证书，准备使用自签证书", args...)
	}

	// 如果没有拿到合适的证书，就返回自签证书。
	if self := cp.self.Load(); self != nil {
		cp.log.Debug("返回自签证书", attrs...)
		return self, nil
	}

	cp.log.Info("准备开始自签证书", attrs...)
	if crt, err := cp.selfSigned(); crt != nil {
		cp.log.Debug("生成并返回自签证书", attrs...)
		return crt, nil
	} else {
		attrs = append(attrs, "error", err)
		cp.log.Warn("自签证书生成错误", attrs...)
		return nil, err
	}
}

func (cp *certPool) Reset() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.cert.Store(nil)
	cp.self.Store(nil)
}

func (cp *certPool) slowLoad() *certMap {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	if crtm := cp.cert.Load(); crtm != nil {
		return crtm
	}

	crtm := &certMap{certs: make(map[string][]*tls.Certificate, 16)}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pairs, err := cp.slow(ctx)
	if err != nil {
		crtm.err = err
		// 可能是网络波动导致的超时问题，不放入结果缓存。
		if te, ok := err.(interface{ Timeout() bool }); !ok || !te.Timeout() {
			cp.cert.Store(crtm)
		}

		return crtm
	}

	for _, kp := range pairs {
		crtm.put(kp)
	}
	cp.cert.Store(crtm)

	return crtm
}

func (cp *certPool) selfSigned() (*tls.Certificate, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	if self := cp.self.Load(); self != nil {
		return self, nil
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, err
	}
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := &x509.Certificate{
		IsCA:                  true,
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{CommonName: "aegis", Organization: []string{"aegis"}},
		NotBefore:             now.AddDate(0, 0, -1),
		NotAfter:              now.AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{{127, 0, 0, 1}},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	cp.self.Store(&tlsCert)

	return &tlsCert, nil
}

type certMap struct {
	err   error
	certs map[string][]*tls.Certificate
}

func (cm *certMap) Match(sni string) (*tls.Certificate, error) {
	if cm.err != nil || len(cm.certs) == 0 {
		return nil, cm.err
	}

	now := time.Now()
	// https://github.com/golang/go/blob/go1.22.5/src/crypto/tls/common.go#L1141-L1154
	name := strings.ToLower(sni)
	var last *tls.Certificate
	for _, crt := range cm.certs[name] {
		notBefore, notAfter := crt.Leaf.NotBefore, crt.Leaf.NotAfter
		if now.After(notBefore) && now.Before(notAfter) {
			return crt, nil
		}
		last = crt
	}

	labels := strings.Split(name, ".")
	labels[0] = "*"
	wildcardName := strings.Join(labels, ".")
	for _, crt := range cm.certs[wildcardName] {
		notBefore, notAfter := crt.Leaf.NotBefore, crt.Leaf.NotAfter
		if now.After(notBefore) && now.Before(notAfter) {
			return crt, nil
		}
		last = crt
	}

	return last, nil
}

func (cm *certMap) put(crt *tls.Certificate) {
	leaf := crt.Leaf
	for _, name := range leaf.DNSNames {
		cm.certs[name] = append(cm.certs[name], crt)
	}
	for _, ip := range leaf.IPAddresses {
		name := ip.String()
		cm.certs[name] = append(cm.certs[name], crt)
	}
}
