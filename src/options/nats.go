package options

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsOptions struct {
	Url        string
	Subject    string
	Tls        bool
	TlsCrtPath string
	TlsKeyPath string
}

func (n *NatsOptions) Connect() (*nats.Conn, error) {
	natsOptions := make([]nats.Option, 0)
	certificates := make([]tls.Certificate, 0)

	if n.Tls {
		if n.TlsCrtPath != "" || n.TlsKeyPath != "" {
			cert, err := tls.LoadX509KeyPair(n.TlsCrtPath, n.TlsKeyPath)
			if err != nil {
				log.Fatalf("error parsing X509 certificate/key pair: %v", err)
			}
			certificates = append(certificates, cert)
		}

		tlsConfig := &tls.Config{
			Certificates: certificates,
			MinVersion:   tls.VersionTLS12,
		}
		natsOptions = append(natsOptions, nats.Secure(tlsConfig))
	}

	return nats.Connect(n.Url, natsOptions...)
}

func NewNatsOptions(prefix string, name string) *NatsOptions {
	natsFlags := NatsOptions{}
	flag.StringVar(&natsFlags.Url, fmt.Sprintf("%s-url", prefix), nats.DefaultURL, fmt.Sprintf("%s url", name))
	flag.StringVar(&natsFlags.Subject, fmt.Sprintf("%s-subj", prefix), "*", fmt.Sprintf("%s url", name))
	return &natsFlags
}
