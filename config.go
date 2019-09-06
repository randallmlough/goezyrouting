package goezyrouting

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
)

// Load returns Configuration struct
func LoadConfig(path string) (*Config, error) {
	ext := filepath.Ext(path)
	if ext != ".json" {
		return nil, errors.New("unsupported file")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}
	var cfg = new(Config)

	if err := json.Unmarshal(bytes, cfg); err != nil {
		return nil, errors.Wrap(err, "unable to decode into struct")
	}

	return cfg, nil
}

type Config struct {
	Application *Application   `json:"application,omitempty"`
	Website     *WebsiteConfig `json:"website,omitempty"`
}

type Application struct {
	Port string `json:"port,omitempty"`
}

type WebsiteConfig struct {
	Title          string     `json:"title,omitempty"`
	Favicon        string     `json:"favicon,omitempty"`
	Assets         FileSystem `json:"assets,omitempty"`
	PublicFolder   FileSystem `json:"public_folder,omitempty"`
	TemplateDir    string     `json:"template_dir,omitempty"`
	TemplatePrefix string     `json:"template_prefix,omitempty"`
	BaseTemplate   string     `json:"base_template,omitempty"`
	ViewsDir       string     `json:"views_dir,omitempty"`
	Styles         []string   `json:"styles,omitempty"`
}

type FileSystem struct {
	Path string `json:"path,omitempty"`
	Dir  string `json:"dir,omitempty"`
}
