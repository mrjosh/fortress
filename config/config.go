package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfMap struct {
	Items map[string]*DomainConfig
}

type DomainConfig struct {
	AccessToken string `yaml:"access_token"`
}

func LoadFromReader(reader io.Reader) (*ConfMap, error) {
	cfg := new(ConfMap)
	if err := yaml.NewDecoder(reader).Decode(&cfg.Items); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Loadfile(filename string) (*ConfMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := LoadFromReader(f)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *ConfMap) GetConfigByDomain(domain string) (*DomainConfig, error) {
	domainCfg, ok := cfg.Items[domain]
	if ok {
		return domainCfg, nil
	}
	return nil, fmt.Errorf("unknown server name [%s]", domain)
}
