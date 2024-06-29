package config

import (
	"github.com/caarlos0/env/v11"
	"log"
	"sync"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_DSN" envDefault:"root:root@tcp(localhost:3306)/goph_keeper?charset=utf8mb4"`
	JWT         struct {
		Secret           string `env:"JWT_SECRET" envDefault:"jwt_secret"`
		ExpiresInSeconds uint32 `env:"JWT_EXPIRES_IN_SECONDS" envDefault:"3600"`
	}
	GRPC struct {
		ServerAddr string `env:"SERVER_ADDR" envDefault:":3200"`
		TLS        struct {
			// CertPEM путь к файлу с сертификатом.
			CertPEM string `env:"CERT_PEM" envDefault:"x509/server_cert.pem"`
			// KeyPEM путь к файлу с приватным ключом.
			KeyPEM string `env:"KEY_PEM" envDefault:"x509/server_key.pem"`
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
