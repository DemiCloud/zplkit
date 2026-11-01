package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// zplBreakers are characters that start a new ZPL command (^) or control
// instruction (~); a config value containing one could inject commands
// into a rendered label rather than just appearing as label text.
const zplBreakers = "^~"

// Load reads and validates the config file at path. A missing file is not
// an error — it returns a zero-value Config, since a template that doesn't
// reference {{.CompanyName}}/{{.AssetURL}} doesn't need one.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	dec := toml.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) validate() error {
	if strings.ContainsAny(c.CompanyName, zplBreakers) {
		return fmt.Errorf("company_name must not contain %q", zplBreakers)
	}
	if strings.ContainsAny(c.AssetURL, zplBreakers) {
		return fmt.Errorf("asset_url must not contain %q", zplBreakers)
	}
	switch c.Printer.Backend {
	case "", "socket", "ipp":
	default:
		return fmt.Errorf("printer.backend must be %q or %q, got %q", "socket", "ipp", c.Printer.Backend)
	}
	return nil
}
