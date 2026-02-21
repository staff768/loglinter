package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
}

func LoadSensitiveKeywords(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg.Keywords, nil
}
