package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfMap struct {
	Domains map[string]*DomainConfig
}

type DomainConfig struct {
	AccessToken string `yaml:"access_token"`
}

func LoadFromReader(reader io.Reader) (*ConfMap, error) {
	cfg := new(ConfMap)
	if err := yaml.NewDecoder(reader).Decode(&cfg.Domains); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadFile(filename string) (*ConfMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &ConfMap{
				Domains: make(map[string]*DomainConfig),
			}, nil
		}
		return nil, err
	}
	cfg, err := LoadFromReader(f)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *ConfMap) GetGitDomainConfig(domain string) (*DomainConfig, bool) {
	domainCfg, ok := cfg.Domains[domain]
	return domainCfg, ok
}
