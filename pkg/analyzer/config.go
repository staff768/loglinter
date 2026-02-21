package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SensitiveDataConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Keywords []string `yaml:"keywords"`
}

type RulesConfig struct {
	Lowercase     bool                `yaml:"lowercase"`
	EnglishOnly   bool                `yaml:"english_only"`
	SpecialChars  bool                `yaml:"special_chars"`
	SensitiveData SensitiveDataConfig `yaml:"sensitive_data"`
}

type Config struct {
	Rules RulesConfig `yaml:"rules"`
}

func LoadConfig(path string) (*Config, error) {
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
	return &cfg, nil
}
