package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

func createProject(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	f, err := common.OpenArg(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var proj oapi.HTCProject
	if err := dec.Decode(&proj); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.CreateProject(ctx, oapi.NewOptHTCProject(proj))
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProject:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.CreateProjectForbidden, *oapi.CreateProjectUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var CreateCmd = &cobra.Command{
	Use:   "create PROJECT_JSON",
	Short: "Creates an HTC project using JSON file or - for stdin.",
	Run:   common.WrapRunE(createProject),
	Args:  cobra.ExactArgs(1),
}

func init() {
	// Prepare a sample JSON payload for our example.
	proj := oapi.HTCProject{
		ProjectDescription: "a description of the project",
		ProjectName:        "a name for the project",
		Regions: []oapi.RescaleRegion{
			oapi.RescaleRegionAZUREUSSOUTHCENTRAL,
		},
	}
	b, err := json.MarshalIndent(&proj, "", "  ")
	if err != nil {
		panic("Unable to serialize `project create` JSON example: " + err.Error())
	}

	CreateCmd.Long = CreateCmd.Short + `

The global list of region names can be found by running:

	htc region get

Note that your workspace likely does *NOT* have access to all regions.
`

	// NB: EOF has no leading space so that copy/paste works w/o
	// editing.
	CreateCmd.Example = fmt.Sprintf(`
htc project create - <<EOF
  %s
EOF`, string(b))

}
