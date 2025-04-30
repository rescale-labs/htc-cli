package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

var htcJobFields = []Field{
	Field{"Job UUID", "%-38s", "%-38s"},
	Field{"Created", "%19s", "%19s"},
	Field{"Completed", "%19s", "%19s"},
	Field{"Status", "%21s", "%21s"},
}

type HTCJob oapi.HTCJob

func (j *HTCJob) Fields() []Field {
	return htcJobFields
}

// Because of how ogen constructs the underlying type, json output
// fails unless we use the concrete type's UnmarshalJSON. So, cast
// back to it and call it instead of letting the stdlib try on its own
// with reflection.
func (j *HTCJob) MarshalJSON() ([]byte, error) {
	c := (*oapi.HTCJob)(j)
	return c.MarshalJSON()
}

func (j *HTCJob) WriteRows(rowFmt string, w io.Writer) error {
	_, err := fmt.Fprintf(
		w, rowFmt,
		j.JobUUID.Value,
		formatDateTime(j.CreatedAt),
		formatDateTime(j.CompletedAt),
		j.Status.Value,
	)
	return err
}

type HTCJobs []oapi.HTCJob

func (s HTCJobs) Fields() []Field {
	return htcJobFields
}

func (s HTCJobs) WriteRows(rowFmt string, w io.Writer) error {
	for _, j := range s {
		if err := (*HTCJob)(&j).WriteRows(rowFmt, w); err != nil {
			return err
		}
	}
	return nil
}

var htcJobStatusEventFields = []Field{
	Field{"Time", "%19s", "%19s"},
	Field{"Status", "%21s", "%21s"},
	Field{"Status Reason", "%19s", "%19s"},
	Field{"Container Exit Code", "%19s", "%19d"},
	Field{"Container Exit Reason", "%24s", "%24s"},
	Field{"Instance Type", "%20s", "%20s"},
	Field{"CSP", "%5s", "%5s"},
	Field{"Region", "%15s", "%15s"},
	Field{"Priority", "%10s", "%10s"},
	Field{"Instance ID", "%20s", "%20s"},
}

type HTCJobStatusEvent oapi.RescaleJobStatusEvent

func (e *HTCJobStatusEvent) Fields() []Field {
	return htcJobStatusEventFields
}

func (e *HTCJobStatusEvent) WriteRows(rowFmt string, w io.Writer) error {
	container := e.Container.Value
	instanceLabels := e.InstanceLabels.Value
	_, err := fmt.Fprintf(
		w, rowFmt,
		formatDateTime(e.DateTime),
		e.Status.Value,
		e.StatusReason.Value,
		int(container.ExitCode.Value),
		container.Reason.Value,
		instanceLabels.InstanceType.Value,
		instanceLabels.Csp.Value,
		instanceLabels.Region.Value,
		instanceLabels.Priority.Value,
		e.InstanceId.Value,
	)
	return err
}

type HTCJobStatusEvents []oapi.RescaleJobStatusEvent

func (s HTCJobStatusEvents) Fields() []Field {
	return htcJobStatusEventFields
}

func (s HTCJobStatusEvents) WriteRows(rowFmt string, w io.Writer) error {
	for _, j := range s {
		if err := (*HTCJobStatusEvent)(&j).WriteRows(rowFmt, w); err != nil {
			return err
		}
	}
	return nil
}
