package worker

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScrapeGroceryTask(t *testing.T) {
	targetURL := "https://www.tokopedia.com/superindo"

	task, err := NewScrapeGroceryTask(targetURL)
	require.NoError(t, err)

	// Assert the task type name is correct
	assert.Equal(t, TaskScrapeGrocery, task.Type())

	// Assert the payload is properly encoded
	var payload ScrapeGroceryPayload
	err = json.Unmarshal(task.Payload(), &payload)
	require.NoError(t, err)
	assert.Equal(t, targetURL, payload.TargetURL)
}

func TestParseScrapeGroceryPayload_Valid(t *testing.T) {
	targetURL := "https://www.superindo.co.id"
	task, err := NewScrapeGroceryTask(targetURL)
	require.NoError(t, err)

	payload, err := ParseScrapeGroceryPayload(task)
	require.NoError(t, err)
	assert.Equal(t, targetURL, payload.TargetURL)
}

func TestParseScrapeGroceryPayload_EmptyURL(t *testing.T) {
	// Manually create a task with an empty URL to simulate bad payload
	task, err := NewScrapeGroceryTask("")
	require.NoError(t, err)

	_, err = ParseScrapeGroceryPayload(task)
	assert.Error(t, err, "expected an error for empty target_url")
	assert.Contains(t, err.Error(), "target_url is required")
}
