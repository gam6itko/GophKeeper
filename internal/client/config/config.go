package config

import (
	"github.com/caarlos0/env/v11"
	"log"
	"sync"
)

type Config struct {
	GRPC struct {
		ServerAddr string `env:"SERVER_ADDR" envDefault:":3200"`
		TLS        struct {
			// CertPEM путь к файлу с сертификатом.
			CertPEM    string `env:"CERT_PEM" envDefault:"x509/ca_cert.pem"`
			ServerHost string `env:"SERVER_NAME" envDefault:"x.test.example.com"`
		}
	}
}

func Load() Config {
	fn := sync.OnceValue(func() Config {
		cfg := Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatal(err)
		}
		return cfg
	})

	return fn()
}
