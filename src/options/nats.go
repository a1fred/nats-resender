package options

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func NewNatsOptions(prefix string, name string) *NatsOptions {
	natsFlags := NatsOptions{Name: name}
	flag.StringVar(&natsFlags.Url, fmt.Sprintf("%s-url", prefix), nats.DefaultURL, fmt.Sprintf("%s url", name))
	flag.StringVar(&natsFlags.Subject, fmt.Sprintf("%s-subj", prefix), "*", fmt.Sprintf("%s url", name))

	flag.StringVar(&natsFlags.UserCreds, fmt.Sprintf("%s-creds", prefix), "", fmt.Sprintf("%s user credentials file", name))
	flag.StringVar(&natsFlags.NkeyFile, fmt.Sprintf("%s-nkey", prefix), "", fmt.Sprintf("%s nkey seed file", name))

	flag.StringVar(&natsFlags.TlsCrtPath, fmt.Sprintf("%s-tlscrt", prefix), "", fmt.Sprintf("%s tls cert file", name))
	flag.StringVar(&natsFlags.TlsKeyPath, fmt.Sprintf("%s-tlskey", prefix), "", fmt.Sprintf("%s tls key file", name))
	flag.StringVar(&natsFlags.TlsCaCrtPath, fmt.Sprintf("%s-tlscacrt", prefix), "", fmt.Sprintf("%s tls ca cert file", name))
	flag.BoolVar(&natsFlags.TlsInsecure, fmt.Sprintf("%s-tlsinsecure", prefix), false, fmt.Sprintf("%s disable tls cert verification", name))
	return &natsFlags
}

type NatsOptions struct {
	Name string

	Url     string
	Subject string

	UserCreds string
	NkeyFile  string

	TlsCrtPath   string
	TlsKeyPath   string
	TlsCaCrtPath string
	TlsInsecure  bool
}

func (n *NatsOptions) Connect() (*nats.Conn, error) {
	opts := make([]nats.Option, 0)

	opts = append(opts, natsErrorHandler(n.Name))

	// https://github.com/nats-io/nats.go/blob/main/examples/nats-sub/main.go

	if n.UserCreds != "" && n.NkeyFile != "" {
		log.Fatal("specify nkey or creds, not both")
	}

	// Use UserCredentials
	if n.UserCreds != "" {
		log.Printf("Using tls CA cert: %s\n", n.TlsCrtPath)
		opts = append(opts, nats.UserCredentials(n.UserCreds))
	}

	// Use Nkey authentication.
	if n.NkeyFile != "" {
		log.Printf("Using nkey file: %s\n", n.TlsCrtPath)
		opt, err := nats.NkeyOptionFromSeed(n.NkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	// Use TLS client authentication
	if n.TlsCrtPath != "" || n.TlsKeyPath != "" {
		if n.TlsCrtPath == "" || n.TlsKeyPath == "" {
			log.Fatal("specify both tls crt and key to use tls client auth")
		}
		log.Printf("Using tls cert: %s\n", n.TlsCrtPath)
		opts = append(opts, nats.ClientCert(n.TlsCrtPath, n.TlsKeyPath))
	}

	// Use specific CA certificate
	if n.TlsCaCrtPath != "" {
		log.Printf("Using tls CA cert: %s\n", n.TlsCrtPath)
		opts = append(opts, nats.RootCAs(n.TlsCaCrtPath))
	}

	if n.TlsInsecure {
		log.Printf("TLS cert validation disabled\n")
		opts = append(opts, natsTlsInsecureSkipVerifyOption(n.TlsInsecure))
	}

	return nats.Connect(n.Url, opts...)
}

// natsTlsInsecureSkipVerifyOption set tls.Config.InsecureSkipVerify
func natsTlsInsecureSkipVerifyOption(isInsecureSkipVerify bool) nats.Option {
	return func(o *nats.Options) error {
		if o.TLSConfig == nil {
			o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		o.TLSConfig.InsecureSkipVerify = isInsecureSkipVerify
		return nil
	}
}

func natsErrorHandler(name string) nats.Option {
	return func(o *nats.Options) error {
		o.AsyncErrorCB = func(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
			switch natsErr {
			case nats.ErrSlowConsumer:
				{
					pendingMsgs, _, err := sub.Pending()
					if err != nil {
						log.Printf("%s couldn't get pending messages: %v", name, err)
						return
					}
					log.Printf("%s falling behind with %d pending messages on subject %q.\n", name, pendingMsgs, sub.Subject)
				}
			default:
				log.Printf("%s error: %v\n", name, natsErr)
			}
		}

		return nil
	}
}
