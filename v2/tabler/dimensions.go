package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

type ComputeEnvs oapi.HTCProjectDimensions

func (c *ComputeEnvs) Fields() []Field {
	return []Field{
		Field{"Region", "%-24s", "%-24.24s"},
		Field{"Machine Type", "%-20s", "%-20.20s"},
		Field{"Priority", "%-8s", "%-8.8s"},
		Field{"HT", "%-2s", " %1.1s"},
		Field{"Arch", "%-5s", "%-5s"},
		Field{"vCPU", "%-4s", "%4d"},
		Field{"Mem", "%-3s", "%-3.0f"},
		Field{"Swap", "%-4s", "%4s"},
		Field{"Scaling", "%-20s", "%-20.20s"},
	}
}

func formatPriority(p oapi.HTCComputeEnvironmentPriority) string {
	switch p {
	case oapi.HTCComputeEnvironmentPriorityONDEMANDPRIORITY:
		return "ODP"
	case oapi.HTCComputeEnvironmentPriorityONDEMANDECONOMY:
		return "ODE"
	}
	return "UNK"
}

func formatBool(v bool) string {
	if v {
		return "T"
	}
	return "F"
}

func (c *ComputeEnvs) WriteRows(rowFmt string, w io.Writer) error {
	for _, e := range []oapi.HTCComputeEnvironment(*c) {
		_, err := fmt.Fprintf(
			w, rowFmt,
			e.Region.Value,
			e.MachineType.Value,
			formatPriority(e.Priority.Value),
			formatBool(e.Hyperthreading.Value),
			e.Derived.Value.Architecture.Value,
			e.Derived.Value.VCPUs.Value,
			e.Derived.Value.Memory.Value,
			formatBool(e.Derived.Value.Swap.Value),
			e.ComputeScalingPolicy.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
