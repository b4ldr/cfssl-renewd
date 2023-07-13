package config

import (
	request "github.com/b4ldr/cfssl-renewd/internal"
)

type Config struct {
	Signer struct {
		Host string `yaml:"host" env:"SIGNER_HOST" env-description:"CFSSL server host" env-default:"foo"`
		Port string `yaml:"port" env:"SIGNER_PORT" env-description:"CFSSL server port"`
		Cert string `yaml:"cert" env:"SIGNER_CERT" env-description:"TLS cert used for mutual auth"`
		Key  string `yaml:"key" env:"SIGNER_PASSWORD" env-description:"TLS private key used for mutual auth"`
		CA   string `yaml:"ca" env:"SIGNER_CA" env-description:"Path to ca bundle used by cfssl server"`
	} `yaml:"signer"`
	Requests []request.Request `yaml:"requests"`
}
