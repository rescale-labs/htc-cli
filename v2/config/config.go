package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
)

type Config struct {
	// Path to config dir
	DirPath string

	// Rescale API key
	ApiKey string

	// Rescale API url
	ApiUrl string

	// Bearer token
	BearerToken BearerToken

	// Output format
	OutputFormat string
}

type BearerToken struct {
	Value      string    `json:"value"`
	ExpiresIn  int64     `json:"expiresIn"`
	ValidUntil time.Time `json:"validUntil"`
}

func (c *Config) tokenPath() string {
	return path.Join(c.DirPath, "bearer.token")
}

func (c *Config) makeDir() error {
	return os.MkdirAll(c.DirPath, 0700)
}

func (c *Config) NeedsToken(now time.Time) bool {
	return c.BearerToken.ValidUntil.Before(now)
}

func (c *Config) SetToken(t *oapi.HTCToken, issuedAt time.Time) {
	c.BearerToken.Value = t.GetTokenValue().Value
	c.BearerToken.ExpiresIn = t.GetExpiresIn().Value
	c.BearerToken.ValidUntil = issuedAt.Add(
		time.Duration(c.BearerToken.ExpiresIn) * time.Second) //.Format(time.RFC3339)
}

func (c *Config) SaveToken() error {
	if err := c.makeDir(); err != nil {
		return fmt.Errorf("SaveToken: %w", err)
	}

	p := c.tokenPath() + ".tmp"
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("SaveToken: %w", err)
	}
	defer f.Close()

	if err := os.Chmod(p, 0600); err != nil {
		return fmt.Errorf("SaveToken: %w", err)
	}

	e := json.NewEncoder(f)
	if err := e.Encode(&c.BearerToken); err != nil {
		return fmt.Errorf("SaveToken: %w", err)
	}
	if err := os.Rename(p, c.tokenPath()); err != nil {
		return fmt.Errorf("SaveToken: %w", err)
	}
	return nil
}

// Attempts to load a token from our config dir, returning error only
// when token is invalid or we weren't able to read the existing file.
func (c *Config) LoadToken() error {
	p := c.tokenPath()
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("LoadToken: %w", err)
	}
	defer f.Close()

	d := json.NewDecoder(f)
	if err := d.Decode(&c.BearerToken); err != nil {
		return fmt.Errorf("LoadToken: %w", err)
	}
	return nil
}

// Implement oapi.SecuritySource. Note that this *always* is sent up
// with Authorization: Bearer, so when we use a token we need to hack
// that out upstream.
func (c *Config) SecurityScheme(ctx context.Context, operationName string) (oapi.SecurityScheme, error) {
	switch operationName {
	case "AuthTokenGet", "AuthWhoamiGet":
		// These two methods only use API key. See
		// common.ClientWrapper.Do() for companion code.
		return oapi.SecurityScheme{Token: c.ApiKey}, nil
	default:
		return oapi.SecurityScheme{Token: c.BearerToken.Value}, nil
	}
}

// Custom error type used to differentiate between runtime errors and
// errors in argument/parameters passed to a command.
//
// Currently here in config, because config is where many common
// parameters are validated.
type UsageError struct {
	error
}

func UsageErrorf(msg string, args ...interface{}) error {
	return &UsageError{fmt.Errorf(msg, args...)}
}

func NewConfig(dirPath, apiKey, apiUrl, outputFormat string) (*Config, error) {
	c := &Config{
		DirPath:      dirPath,
		ApiKey:       apiKey,
		ApiUrl:       apiUrl,
		OutputFormat: outputFormat,
	}
	if c.DirPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		c.DirPath = path.Join(configDir, "rescale", "htc")
	}

	if c.ApiKey == "" {
		c.ApiKey = os.Getenv("RESCALE_API_KEY")
	}

	if c.ApiUrl == "" {
		c.ApiUrl = "https://htc.rescale.com/api/v1/"
	}
	if _, err := url.Parse(c.ApiUrl); err != nil {
		return nil, UsageErrorf("Invalid API URL: %w", err)
	}

	if c.OutputFormat == "" {
		c.OutputFormat = "json"
	}
	switch c.OutputFormat {
	case "json", "yaml":
		break
	default:
		return nil, UsageErrorf(
			`Output format must be "json" or "yaml" not %q.`,
			c.OutputFormat)
	}

	return c, nil
}
