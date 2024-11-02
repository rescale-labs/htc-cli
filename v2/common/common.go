package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-faster/yaml"
	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/config"
)

type ClientWrapper struct {
	*http.Client
}

func (c *ClientWrapper) Do(r *http.Request) (*http.Response, error) {
	// GET /auth/token must use Token, not Bearer, in its auth.
	if r.Method == "GET" {
		switch {
		case strings.HasSuffix(r.URL.Path, "/auth/token"),
			strings.HasSuffix(r.URL.Path, "/auth/whoami"):
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				r.Header.Set("Authorization", "Token "+auth[7:])
			}
		}
	}
	// res, err := c.Client.Do(r)
	return c.Client.Do(r)
}

func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, config.UsageErrorf("Error setting output format: %w", err)
	}

	config, err := config.NewConfig(
		"",
		"",
		"",
		outputFormat)
	if err != nil {
		return nil, err
	}
	if err := config.LoadToken(); err != nil {
		return nil, err
	}
	return config, nil
}

func getBearerToken(c *oapi.Client) (*oapi.HTCToken, error) {
	ctx := context.Background()
	res, err := c.AuthTokenGet(ctx)
	if err != nil {
		log.Fatalf("Login failed: %s", err)
	}
	switch r := res.(type) {
	case *oapi.HTCToken:
		return r, nil

	case *oapi.OAuth2ErrorResponse:
		return nil, fmt.Errorf("Login failed: error=%s description=%s",
			r.GetError().Value, r.GetErrorDescription().Value)
	}
	return nil, fmt.Errorf("Login failed: unknown response.")
}

func updateBearerToken(c *oapi.Client, config *config.Config) error {
	start := time.Now()
	t, err := getBearerToken(c)
	if err != nil {
		return fmt.Errorf("updateBearerToken: %w", err)
	}
	log.Printf("Bearer token: ExpiresIn=%d", t.GetExpiresIn().Value)
	config.SetToken(t, start)
	return nil
}

type Runner struct {
	Client  *oapi.Client
	Command *cobra.Command
	Config  *config.Config
}

func NewRunner(cmd *cobra.Command) (*Runner, error) {
	cnf, err := loadConfig(cmd)
	if err != nil {
		if _, ok := err.(*config.UsageError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("Configuration failed: %w", err)
	}

	c, err := oapi.NewClient(cnf.ApiUrl, cnf,
		oapi.WithClient(&ClientWrapper{http.DefaultClient}))
	if err != nil {
		return nil, fmt.Errorf("API client initialization failed: %w", err)
	}
	return &Runner{
		Client:  c,
		Command: cmd,
		Config:  cnf,
	}, nil
}

func (r *Runner) UpdateToken(now time.Time) error {
	if r.Config.NeedsToken(time.Now()) {
		if r.Config.ApiKey == "" {
			return fmt.Errorf("Needed a current bearer token, but unable to get one. Please set RESCALE_API_KEY and retry.")
		}
		return r.RenewToken()
	}
	return nil
}

func (r *Runner) RenewToken() error {
	if err := updateBearerToken(r.Client, r.Config); err != nil {
		return fmt.Errorf("API client auth failed: %w", err)
	}
	if err := r.Config.SaveToken(); err != nil {
		return fmt.Errorf("Saving bearer token failed: %s", err)
	}
	return nil
}

// Returns a new runner with an up to date JWT token.
// Use this for any command except those related to auth.
func NewRunnerWithToken(cmd *cobra.Command, now time.Time) (*Runner, error) {
	runner, err := NewRunner(cmd)
	if err != nil {
		return nil, err
	}
	if err := runner.UpdateToken(now); err != nil {
		return nil, err
	}
	return runner, nil
}

type IDParams struct {
	ProjectId string
	TaskId    string
	JobId     string

	RequireProjectId bool
	RequireTaskId    bool
	RequireJobId     bool
}

func (r *Runner) GetIds(p *IDParams) error {
	var errors []error
	if p.RequireProjectId {
		projectId, err := r.Command.Flags().GetString("project-id")
		if err != nil {
			errors = append(errors,
				config.UsageErrorf("Error setting project ID: %w", err))
		} else if projectId == "" {
			errors = append(errors,
				config.UsageErrorf("Error: missing project ID."))
		} else {
			p.ProjectId = projectId
		}
	}

	if p.RequireTaskId {
		taskId, err := r.Command.Flags().GetString("task-id")
		if err != nil {
			errors = append(errors,
				config.UsageErrorf("Error setting task ID: %w", err))
		} else if taskId == "" {
			errors = append(errors,
				config.UsageErrorf("Error: missing task ID."))
		} else {
			p.TaskId = taskId
		}
	}

	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 0 {
		var words []string
		var args []interface{}
		for _, err := range errors {
			words = append(words, "%w")
			args = append(args, err)
		}
		msg := "Errors:\n" + strings.Join(words, "\n  * ")
		return config.UsageErrorf(msg, args...)
	}

	return nil
}

func (r *Runner) PrintResult(res interface{}, w io.Writer) error {
	switch r.Config.OutputFormat {
	case "yaml":
		// NB: YAML encoding doesn't work properly for ogen's OptString,
		// OptInt, etc. types, and there's not an easy way to fix that
		// without us adding code to the same package as what we
		// generate. (Though, I guess that'd be OK to do. We probably
		// have to do the same eventually anyway for plain text output
		// of various entities.)
		e := yaml.NewEncoder(w)
		defer e.Close()
		return e.Encode(res)

	default:
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		return e.Encode(res)
	}
}

func WrapRunE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := f(cmd, args)
		if err != nil {
			if _, ok := err.(*config.UsageError); ok {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n", err)
				cmd.Usage()
				os.Exit(1)
			}
		}
		cobra.CheckErr(err)
	}
}
