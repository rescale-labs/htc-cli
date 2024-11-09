package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	g := GlobalConf{
		ApiUrl: "https://htc.rescale.com/api/v1/",
	}

	cfgPath := filepath.Join(t.TempDir(), "config.toml")
	if err := writeConfig(cfgPath, &g); err != nil {
		t.Fatalf("writeConfig failed: %s", err)
		return
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		t.Fatalf("failed to open %s: %s", cfgPath, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read %s: %s", cfgPath, err)
	}

	got := string(b)
	want := fmt.Sprintf(
		`# Rescale HTC API URL
api_url = '%s'
# Default for --output-format option
output_format = '%s'
# Name of currently selected context. Can be overridden by --context
selected_context = ''
`,
		g.ApiUrl,
		g.OutputFormat,
	)
	if want != got {
		t.Logf("wanted: %q", want)
		t.Logf("got:    %q", got)
		t.Errorf("wanted:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestReadConfig(t *testing.T) {
	want := GlobalConf{
		ApiUrl: "https://htc.rescale.com/api/v1/",
	}

	// TODO: remove duplication w/TestWriteConfig
	b := []byte(fmt.Sprintf(
		`api_url = '%s'
output_format = '%s'
`,
		want.ApiUrl,
		want.OutputFormat,
	))

	cfgPath := filepath.Join(t.TempDir(), "config.toml")
	f, err := os.Create(cfgPath)
	if err != nil {
		t.Fatalf("failed to create %s: %s", cfgPath, err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		t.Fatalf("failed to write %s: %s", cfgPath, err)
	}

	got := GlobalConf{}
	if err := readConfig(cfgPath, &got); err != nil {
		t.Fatalf("readConfig failed: %s", err)
		return
	}

	if !reflect.DeepEqual(want, got) {
		wantJSON, err := json.MarshalIndent(want, "  ", "  ")
		if err != nil {
			t.Fatalf("Unable to serialize %v: %s", want, err)
		}
		gotJSON, err := json.MarshalIndent(got, "  ", "  ")
		if err != nil {
			t.Fatalf("Unable to serialize %v: %s", got, err)
		}
		t.Errorf("want:\n\n%s\n\ngot:\n\n%s", wantJSON, gotJSON)
	}
}
