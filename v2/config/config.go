package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/spf13/cobra"
)

const defaultApiUrl = "https://htc.rescale.com/api/v1/"
const DefaultContextName = "default"

// All the configuration for common commands goes here and will be set
// on initialization of config. If you have a flag that's used by
// multiple commands (e.g. --output, --project-id, --task-id), it
// probably belongs here and should be set usings flags, with fallback
// to the config file, during NewConfig().
//
// Aside: --limit is not *yet* in here, but probably also belongs here.
type Config struct {
	// Path to config dir
	DirPath string

	// Currently selected config context name
	Context string

	// Rescale API url
	ApiUrl string

	// Active API key and bearer token
	Credentials Credentials

	// Output format
	OutputFormat string

	// Project ID (optional)
	ProjectId string

	// Task ID (optional)
	TaskId string
}

func (c *Config) configPath() string {
	return filepath.Join(c.DirPath, "config.toml")
}

func (c *Config) tokenPath(cfgContext string) string {
	return filepath.Join(c.DirPath, "private", cfgContext+"-credentials.json")
}

// Makes all directories needed for storing config and credentials.
func (c *Config) makeDirs() error {
	if err := os.MkdirAll(c.DirPath, 0755); err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(c.DirPath, "private"), 0700)
}

func (c *Config) NeedsToken(now time.Time) bool {
	return c.Credentials.BearerToken.ValidUntil.Before(now)
}

func (c *Config) SetToken(t *oapi.HTCToken, issuedAt time.Time) {
	c.Credentials.BearerToken = BearerToken{
		Value:     t.GetTokenValue().Value,
		ExpiresIn: t.GetExpiresIn().Value,
		ValidUntil: issuedAt.Add(
			time.Duration(t.GetExpiresIn().Value) * time.Second), //.Format(time.RFC3339)
	}
}

func (c *Config) SaveCredentials() error {
	if err := c.makeDirs(); err != nil {
		return fmt.Errorf("SaveCredentials: %w", err)
	}
	if err := writeCredentials(c.tokenPath(c.Context), &c.Credentials); err != nil {
		return fmt.Errorf("SaveCredentials: %w", err)
	}
	return nil
}

// Attempts to load a token from our config dir, returning error only
// when token can't be parsed or we weren't able to read the existing
// file.
func (c *Config) loadCredentials() error {
	if err := readCredentials(c.tokenPath(c.Context), &c.Credentials); err != nil {
		return fmt.Errorf("Failed to load credentials: %w", err)
	}
	return nil
}

// Implements oapi.SecuritySource. Note that ogen *always* sends this
// token up with "Authorization: Bearer", so when we need to send
// "Authorization: Token" that's hacked out upstream by the custom HTTP
// client in common.ClientWrapper.
func (c *Config) SecurityScheme(ctx context.Context, operationName string) (oapi.SecurityScheme, error) {
	switch operationName {
	case "AuthTokenGet", "AuthWhoamiGet":
		// Only these two methods use the API key. See
		// common.ClientWrapper.Do() for companion code.
		return oapi.SecurityScheme{Token: c.Credentials.ApiKey}, nil
	default:
		return oapi.SecurityScheme{Token: c.Credentials.BearerToken.Value}, nil
	}
}

// Reads global config from c.configPath() and returns a pointer to it,
// or any error.
//
// g.Contexts() will always be initialized to a map.
func (c *Config) ReadGlobalConf() (*GlobalConf, error) {
	var g GlobalConf
	if err := readConfig(c.configPath(), &g); err != nil {
		return nil, err
	}
	if g.Contexts == nil {
		g.Contexts = make(map[string]*ContextConf)
	}
	return &g, nil
}

// Sets a configuration value
func (c *Config) Set(key, value string, global bool) error {
	g, err := c.ReadGlobalConf()
	if err != nil {
		// by this point in the runtime, a config should *always* be
		// present on disk.
		return err
	}

	var fieldName string
	var field reflect.Value

	if !global {
		conf := g.Contexts[c.Context]
		for _, f := range reflect.VisibleFields(reflect.TypeOf(conf).Elem()) {
			toml := strings.SplitN(f.Tag.Get("toml"), ",", 2)[0]
			if key == toml {
				fieldName = f.Name
				if conf == nil {
					conf = &ContextConf{}
					g.Contexts[c.Context] = conf
				}
				// Pass the pointer to this struct in so that reflect can modify
				// it in place.
				field = reflect.ValueOf(conf).Elem().FieldByName(fieldName)
				break
			}
		}
	} else {
		t := reflect.TypeOf(g).Elem()
		for _, f := range reflect.VisibleFields(t) {
			toml := strings.SplitN(f.Tag.Get("toml"), ",", 2)[0]
			if key == toml {
				fieldName = f.Name
				field = reflect.ValueOf(g).Elem().FieldByName(fieldName)
				break
			}
		}

		// Initial config block for the user to edit if we're switching
		// contexts and it doesn't exist yet.
		if key == "selected_context" && g.Contexts[value] == nil {
			g.Contexts[value] = &ContextConf{}
		}
	}

	if fieldName == "" {
		return UsageErrorf(
			`Unknown config key %q. Please review:

    %s

for the list of global and context specific config keys.`,
			key, c.configPath())
	}

	v := reflect.ValueOf(value)
	if v.CanConvert(field.Type()) {
		field.Set(v.Convert(field.Type()))
	} else {
		return UsageErrorf("Unable to convert value %q to type %s", v, field.Type())
	}

	return writeConfig(c.configPath(), g)
}

func (c *Config) Delete(contextName string) error {
	g, err := c.ReadGlobalConf()
	if err != nil {
		// by this point in the runtime, a config should *always* be
		// present on disk.
		return err
	}
	delete(g.Contexts, contextName)
	return writeConfig(c.configPath(), g)
}

// Takes a *cobra.Command and returns the Config needed for running
// it.
//
// Key steps:
//
//  1. Identify the context to use using `--context X` or whatever is
//     defined in ~/.config/rescale/htc/config.toml.
//  2. Populate the Config struct using whatever TOML config has been
//     written for this context.
//  3. Override this config with anything passed in via cmd.Flags().
//  4. Load the credentials for this context.
func NewConfig(cmd *cobra.Command) (*Config, error) {
	c := &Config{}

	// Set dirPath. Currently, this is always from os.UserConfigDir().
	if c.DirPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		c.DirPath = path.Join(configDir, "rescale", "htc")
	}

	// Load global config file if exists.
	var configNotExist bool
	g, err := c.ReadGlobalConf()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			g = &GlobalConf{
				Contexts: make(map[string]*ContextConf),
			}
			configNotExist = true
		} else {
			return nil, fmt.Errorf("Error reading config file: %w", err)
		}
	}

	// Determine current context to use
	flagset := cmd.Flags()
	optContext, err := flagset.GetString("context")
	if err != nil {
		return nil, UsageErrorf("Error reading -context: %w", err)
	}
	switch {
	case optContext != "":
		c.Context = optContext

	case g.SelectedContext != "":
		c.Context = g.SelectedContext

	default:
		c.Context = DefaultContextName
	}

	// Create currently selected context if it does exist.
	contextConf := g.Contexts[c.Context]
	if contextConf == nil {
		contextConf = &ContextConf{}
		g.Contexts[c.Context] = contextConf
	}

	//
	// Apply all other global options
	//

	// --output
	optOutputFormat, err := flagset.GetString("output")
	if err != nil {
		return nil, UsageErrorf("Error reading -output: %w", err)
	}
	switch {
	case optOutputFormat != "":
		c.OutputFormat = optOutputFormat

	case g.OutputFormat != "":
		c.OutputFormat = g.OutputFormat
	}
	switch c.OutputFormat {
	case "json", "yaml", "text":
		break
	default:
		return nil, UsageErrorf(
			`Output format must be one of "json", "yaml", or "text", not %q.`,
			c.OutputFormat)
	}

	// --project-id
	//
	if flagset.Lookup("project-id") != nil {
		projectId, err := cmd.Flags().GetString("project-id")
		if err != nil {
			return nil, UsageErrorf("Error reading --project-id: %w", err)
		}
		if projectId != "" {
			c.ProjectId = projectId
		}
	}
	if c.ProjectId == "" {
		c.ProjectId = contextConf.ProjectId
	}

	// --task-id
	if flagset.Lookup("task-id") != nil {
		taskId, err := cmd.Flags().GetString("task-id")
		if err != nil {
			return nil, UsageErrorf("Error reading --task-id: %w", err)
		}
		if taskId != "" {
			c.TaskId = taskId
		}
	}
	if c.TaskId == "" {
		c.TaskId = contextConf.TaskId
	}

	//
	// Set credentials and remaining config values
	//
	if err := c.loadCredentials(); err != nil {
		return nil, err
	}
	switch {
	case os.Getenv("RESCALE_API_KEY") != "":
		c.Credentials.ApiKey = os.Getenv("RESCALE_API_KEY")
	case c.Credentials.ApiKey == "" && c.NeedsToken(time.Now()):
		return nil, UsageErrorf(
			"API key is not present in %s or as env var RESCALE_API_KEY, "+
				"and bearer token in %s does not exist or is expired.",
			c.tokenPath(c.Context), c.tokenPath(c.Context))
	}

	switch {
	case contextConf.ApiUrl != "":
		c.ApiUrl = contextConf.ApiUrl
	case g.ApiUrl != "":
		c.ApiUrl = g.ApiUrl
	default:
		c.ApiUrl = defaultApiUrl
	}
	if _, err := url.Parse(c.ApiUrl); err != nil {
		return nil, UsageErrorf("Invalid API URL: %w", err)
	}

	// Write config if it did not exist.
	if configNotExist {
		if err := c.makeDirs(); err != nil {
			return nil, err
		}
		// Populate a few defaults
		g.SelectedContext = c.Context
		g.ApiUrl = c.ApiUrl
		g.Contexts[c.Context] = contextConf
		if err := writeConfig(c.configPath(), g); err != nil {
			return nil, err
		}
	}

	return c, nil
}
