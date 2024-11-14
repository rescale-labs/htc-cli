package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

//
// Config
//

type ContextConf struct {
	ApiUrl    string `json:"api_url,omitempty" toml:"api_url" comment:"Rescale HTC API URL"`
	ProjectId string `json:"project_id,omitempty" toml:"project_id" comment:"Default HTC Project Id for commands that take -project-id"`
	TaskId    string `json:"task_id,omitempty" toml:"task_id" comment:"Default HTC Task Id for commands that take -task-id"`
}

type GlobalConf struct {
	ApiUrl          string                  `json:"api_url,omitempty" toml:"api_url" comment:"Rescale HTC API URL"`
	OutputFormat    string                  `json:"output_format,omitempty" toml:"output_format" comment:"Default for --output-format option"`
	SelectedContext string                  `json:"selected_context,omitempty" toml:"selected_context" comment:"Name of currently selected context. Can be overridden by --context"`
	Contexts        map[string]*ContextConf `json:"contexts,omitempty" toml:"contexts" comment:"Context specific configurations"`
}

func writeConfig(cfgPath string, g *GlobalConf) error {
	f, err := os.CreateTemp(filepath.Dir(cfgPath), filepath.Base(cfgPath)+".*")
	if err != nil {
		return err
	}
	defer f.Close()
	e := toml.NewEncoder(f)
	if err := e.Encode(g); err != nil {
		defer os.Remove(f.Name())
		return err
	}
	return os.Rename(f.Name(), cfgPath)
}

func readConfig(cfgPath string, g *GlobalConf) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	d := toml.NewDecoder(f)
	d.DisallowUnknownFields()
	return d.Decode(g)
}

//
// Credentials
//

type BearerToken struct {
	Value      string    `json:"value"`
	ExpiresIn  int64     `json:"expiresIn"`
	ValidUntil time.Time `json:"validUntil"`
}

type Credentials struct {
	ApiKey      string      `json:"api_key"`
	BearerToken BearerToken `json:"bearer_token,omitempty"`
}

func writeCredentials(p string, c *Credentials) error {
	f, err := os.CreateTemp(filepath.Dir(p), filepath.Base(p)+".*")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := os.Chmod(f.Name(), 0600); err != nil {
		return err
	}

	e := json.NewEncoder(f)
	e.SetIndent("", "  ") // be friendly to humans; indent
	if err := e.Encode(c); err != nil {
		return err
	}
	if err := os.Rename(f.Name(), p); err != nil {
		return err
	}
	return nil
}

// Attempts to load a token from our config dir, returning error only
// when token can't be parsed or we weren't able to read the existing
// file.
func readCredentials(p string, c *Credentials) error {
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	return nil
}
