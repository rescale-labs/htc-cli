// Code generated by ogen, DO NOT EDIT.

package _oas

import (
	"bytes"
	"net/http"

	"github.com/go-faster/jx"

	ht "github.com/ogen-go/ogen/http"
)

func encodeCreateTaskRequest(
	req OptHTCTask,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsPostRequest(
	req OptHTCProject,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdDimensionsPutRequest(
	req []HTCComputeEnvironment,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		if req != nil {
			e.ArrStart()
			for _, elem := range req {
				elem.Encode(e)
			}
			e.ArrEnd()
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdLimitsIDPatchRequest(
	req OptHTCLimitUpdate,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdLimitsPostRequest(
	req OptHTCLimitCreate,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdPatchRequest(
	req OptHTCProjectUpdate,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdTaskRetentionPolicyPutRequest(
	req OptTaskRetentionPolicy,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcProjectsProjectIdTasksTaskIdPatchRequest(
	req OptHTCTaskUpdate,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeHtcWorkspacesWorkspaceIdTaskRetentionPolicyPutRequest(
	req OptWorkspaceTaskRetentionPolicy,
	r *http.Request,
) error {
	const contentType = "application/json"
	if !req.Set {
		// Keep request with empty body if value is not set.
		return nil
	}
	e := new(jx.Encoder)
	{
		if req.Set {
			req.Encode(e)
		}
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}

func encodeSubmitJobsRequest(
	req []HTCJobSubmitRequest,
	r *http.Request,
) error {
	const contentType = "application/json"
	e := new(jx.Encoder)
	{
		e.ArrStart()
		for _, elem := range req {
			elem.Encode(e)
		}
		e.ArrEnd()
	}
	encoded := e.Bytes()
	ht.SetBody(r, bytes.NewReader(encoded), contentType)
	return nil
}
