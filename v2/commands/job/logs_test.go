package job

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

// mockJobClient implements the oapi.JobInvoker interface for testing
type mockJobClient struct {
	responses []mockResponse
	callCount int
}

type mockResponse struct {
	logs *oapi.HTCJobLogs
	err  error
}

// we need to implement all JobInvoker interface methods even though only GetLogs is used in our tests
func (m *mockJobClient) CancelJobs(_ context.Context, _ oapi.CancelJobsParams) (oapi.CancelJobsRes, error) {
	return nil, nil
}

func (m *mockJobClient) GetEvents(_ context.Context, _ oapi.GetEventsParams) (oapi.GetEventsRes, error) {
	return nil, nil
}

func (m *mockJobClient) GetJob(_ context.Context, _ oapi.GetJobParams) (oapi.GetJobRes, error) {
	return nil, nil
}

func (m *mockJobClient) GetJobs(_ context.Context, _ oapi.GetJobsParams) (oapi.GetJobsRes, error) {
	return nil, nil
}

func (m *mockJobClient) GetLogs(_ context.Context, _ oapi.GetLogsParams) (oapi.GetLogsRes, error) {
	if m.callCount >= len(m.responses) {
		return nil, nil
	}

	response := m.responses[m.callCount]
	m.callCount++

	if response.err != nil {
		return nil, response.err
	}

	return response.logs, nil
}

func (m *mockJobClient) SubmitJobs(_ context.Context, _ []oapi.HTCJobSubmitRequest, _ oapi.SubmitJobsParams) (oapi.SubmitJobsRes, error) {
	return nil, nil
}

func createLogEvent(timestamp time.Time, message string) oapi.HTCLogEvent {
	return oapi.HTCLogEvent{
		Timestamp: oapi.NewOptInstant(oapi.Instant(timestamp)),
		Message:   oapi.NewOptString(message),
	}
}

func createLogsResponse(events []oapi.HTCLogEvent, nextToken string) *oapi.HTCJobLogs {
	var next oapi.OptURI
	if nextToken != "" {
		u, _ := url.Parse("http://example.com?pageIndex=" + nextToken)
		next = oapi.NewOptURI(*u)
	}

	return &oapi.HTCJobLogs{
		Items: events,
		Next:  next,
	}
}

func TestLogs_HappyPath(t *testing.T) {
	now := time.Now()

	// Create mock responses with log events
	events1 := []oapi.HTCLogEvent{
		createLogEvent(now.Add(-2*time.Second), "First log message"),
		createLogEvent(now.Add(-1*time.Second), "Second log message"),
	}
	events2 := []oapi.HTCLogEvent{
		createLogEvent(now, "Third log message"),
	}

	mockClient := &mockJobClient{
		responses: []mockResponse{
			{logs: createLogsResponse(events1, "page2"), err: nil},
			{logs: createLogsResponse(events2, ""), err: nil}, // No next token = end
		},
	}

	// this should return the two lines
	result1, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result1.Items) != 2 {
		t.Errorf("Expected 2 items in first response, got %d", len(result1.Items))
	}

	result2, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "page2")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result2.Items) != 1 {
		t.Errorf("Expected 1 item in second response, got %d", len(result2.Items))
	}
}

func TestLogs_EmptyResponse(t *testing.T) {
	// we need to verify a real empty response
	mockClient := &mockJobClient{
		responses: []mockResponse{
			{logs: createLogsResponse([]oapi.HTCLogEvent{}, ""), err: nil},
		},
	}

	// should return nothing
	result, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("Expected 0 items in empty response, got %d", len(result.Items))
	}
}

func TestLogs_AwsBugEmptyPagesWithToken(t *testing.T) {
	now := time.Now()

	// it's possible for us to get "fake" empty responses
	// e.g., three pages where the first or second page are empty
	// see https://repost.aws/questions/QUw-BA0Q6PS06SmauQ5rzlSQ/cloudwatch-log-events-not-being-returned-in-cli-and-api
	events := []oapi.HTCLogEvent{
		createLogEvent(now, "Final log message"),
	}

	mockClient := &mockJobClient{
		responses: []mockResponse{
			// fake empty (AWS bug)
			{logs: createLogsResponse([]oapi.HTCLogEvent{}, "page2"), err: nil},
			// another fake empty (AWS bug continues)
			{logs: createLogsResponse([]oapi.HTCLogEvent{}, "page3"), err: nil},
			// real page with logs that we expect to be displayed
			{logs: createLogsResponse(events, ""), err: nil},
		},
	}

	// need to make multiple calls to simulate pagination through empty pages
	result1, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "")
	if err != nil {
		t.Fatalf("Expected no error on first call, got: %v", err)
	}

	if len(result1.Items) != 0 {
		t.Errorf("Expected 0 items in first response (empty page), got %d", len(result1.Items))
	}

	// verify we got a next token despite empty results
	nextToken := result1.Next.Value.Query().Get("pageIndex")
	if nextToken != "page2" {
		t.Errorf("Expected next token 'page2', got '%s'", nextToken)
	}

	result2, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "page2")
	if err != nil {
		t.Fatalf("Expected no error on second call, got: %v", err)
	}

	if len(result2.Items) != 0 {
		t.Errorf("Expected 0 items in second response (empty page), got %d", len(result2.Items))
	}

	result3, err := logs(context.Background(), mockClient, "project-id", "task-id", "job-id", "page3")
	if err != nil {
		t.Fatalf("Expected no error on third call, got: %v", err)
	}

	if len(result3.Items) != 1 {
		t.Errorf("Expected 1 item in third response, got %d", len(result3.Items))
	}
}

func TestWriteRows_EmptyRows(t *testing.T) {
	var output strings.Builder

	err := writeRows([]oapi.HTCLogEvent{}, &output, time.Time{})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output.String() != "" {
		t.Errorf("Expected empty output for empty rows, got: %s", output.String())
	}
}

func TestWriteRows_WithIgnoreBefore(t *testing.T) {
	now := time.Now()
	cutoff := now.Add(-30 * time.Second)

	events := []oapi.HTCLogEvent{
		createLogEvent(now.Add(-60*time.Second), "Old message"), // Before cutoff
		createLogEvent(now.Add(-15*time.Second), "New message"), // After cutoff
	}

	var output strings.Builder
	err := writeRows(events, &output, cutoff)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should only contain the new message
	outputStr := output.String()
	if !strings.Contains(outputStr, "New message") {
		t.Errorf("Expected output to contain 'New message', got: %s", outputStr)
	}
	if strings.Contains(outputStr, "Old message") {
		t.Errorf("Expected output to not contain 'Old message', got: %s", outputStr)
	}
}
