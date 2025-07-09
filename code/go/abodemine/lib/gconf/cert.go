package gconf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Certificate struct {
	CertFile string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"`
	CertBody []byte `json:"cert_body,omitempty" yaml:"cert_body,omitempty"`
	KeyFile  string `json:"key_file,omitempty" yaml:"key_file,omitempty"`
	KeyBody  []byte `json:"key_body,omitempty" yaml:"key_body,omitempty"`
	ClientCa bool   `json:"client_ca,omitempty" yaml:"client_ca,omitempty"`
	RootCa   bool   `json:"root_ca,omitempty" yaml:"root_ca,omitempty"`
}

type Tls struct {
	ClientAuthType tls.ClientAuthType `json:"client_auth_type,omitempty" yaml:"client_auth_type,omitempty"`
	Certificates   []*Certificate     `json:"certificates,omitempty" yaml:"certificates,omitempty"`
}

func GetTLS(config *Tls) (*tls.Config, error) {
	x509ClientCertPool := x509.NewCertPool()
	x509RootCertPool := x509.NewCertPool()

	if len(config.Certificates) == 0 {
		return nil, errors.New("empty conf.Certificates")
	}

	var certs []tls.Certificate

	for i := range config.Certificates {
		c := config.Certificates[i]

		var err error
		var cb, kb []byte

		switch {
		case c.CertFile != "":
			cb, err = os.ReadFile(filepath.Clean(os.ExpandEnv(c.CertFile)))
			if err != nil {
				return nil, fmt.Errorf("failed to read certificates cert file: %w", err)
			}
		case len(c.CertBody) > 0:
			cb = []byte(c.CertBody)
		default:
			return nil, fmt.Errorf("conf.Certificates[%d] must define CertFile or CertBody", i)
		}

		if len(cb) == 0 {
			return nil, fmt.Errorf("conf.Certificates[%d] pem data is empty", i)
		}

		if c.ClientCa {
			ok := x509ClientCertPool.AppendCertsFromPEM(cb)
			if !ok {
				return nil, fmt.Errorf("conf.Pool[%d] failed to be appended to clientCertPool", i)
			}
		}

		if c.RootCa {
			ok := x509RootCertPool.AppendCertsFromPEM(cb)
			if !ok {
				return nil, fmt.Errorf("conf.Pool[%d] failed to be appended to rootCertPool", i)
			}
		}

		if c.ClientCa || c.RootCa {
			// No need to check for keys.
			continue
		}

		switch {
		case c.KeyFile != "":
			kb, err = os.ReadFile(filepath.Clean(os.ExpandEnv(c.KeyFile)))
			if err != nil {
				return nil, fmt.Errorf("failed to read certificates key file: %w", err)
			}
		case len(c.KeyBody) > 0:
			kb = []byte(c.KeyBody)
		default:
			return nil, fmt.Errorf("conf.Certificates[%d] must define KeyFile or KeyBody", i)
		}

		if len(kb) == 0 {
			return nil, fmt.Errorf("conf.Certificates[%d] key data is empty", i)
		}

		tc, err := tls.X509KeyPair(cb, kb)
		if err != nil {
			return nil, fmt.Errorf("failed to create X509KeyPair: %w", err)
		}

		certs = append(certs, tc)
	}

	tc := &tls.Config{
		Certificates: certs,
		ClientAuth:   config.ClientAuthType,
	}

	return tc, nil
}
