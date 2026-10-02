package preview

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gowatchtower/tplprev/internal/notify"
	"github.com/gowatchtower/tplprev/internal/templates"
)

func TestRenderDefaultReport(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["default"],
		[]State{ScannedState, UpdatedState, FailedState},
		nil,
	)
	require.NoError(t, err)
	assert.Contains(t, result, "1 Scanned, 1 Updated")
	assert.Contains(t, result, "1 Failed")
	assert.Contains(t, result, "datamatrix")
	assert.Contains(t, result, "updated to")
}

func TestRenderDefaultReportAppendsRefAndChangelog(t *testing.T) {
	t.Parallel()

	result, err := Render(templates.Templates["default"], []State{UpdatedState}, nil)
	require.NoError(t, err)
	assert.Contains(t, result, "ref "+previewGitRef)
	assert.Contains(t, result, previewChangelog)
}

func TestRenderDefaultReportOmitsRefForFreshContainer(t *testing.T) {
	t.Parallel()

	result, err := Render(templates.Templates["default"], []State{FreshState}, nil)
	require.NoError(t, err)
	assert.Contains(t, result, "fresh")
	assert.NotContains(t, result, "ref ")
}

func TestRenderDefaultLegacyGitMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		give   notify.Entry
		want   string
		absent string
	}{
		{
			name: "changelog line",
			give: notify.Entry{
				Message: "Changelog",
				Data:    map[string]any{"changelog": "https://github.com/org/app/releases"},
			},
			want: "Changelog: https://github.com/org/app/releases",
		},
		{
			name: "found new revision",
			give: notify.Entry{
				Message: "Found new Git revision",
				Data: map[string]any{
					"revision":     "9f2c1ab4d7e5f6a8",
					"short_commit": "9f2c1ab4d7e5",
				},
			},
			want: "Found new Git revision: 9f2c1ab4d7e5f6a8 (9f2c1ab4d7e5)",
		},
		{
			name: "built image from git url context",
			give: notify.Entry{
				Message: "Built image from Git URL context",
				Data: map[string]any{
					"image":    "org/app:git-9f2c1ab4d7e5",
					"image_id": "abc123def456",
				},
			},
			want: "Built image: org/app:git-9f2c1ab4d7e5 (abc123def456)",
		},
		{
			name: "built compose project",
			give: notify.Entry{
				Message: "Built Compose project",
				Data:    map[string]any{"project": "webstack", "service": "api"},
			},
			want: "Built Compose project: webstack (api)",
		},
		{
			name: "git build failed",
			give: notify.Entry{
				Message: "Git build failed. Leaving running container untouched",
				Data: map[string]any{
					"container": "api",
					"error":     "docker build failed",
				},
			},
			want: "Git update failed for api. The running container was left in place: docker build failed",
		},
		{
			name: "compose apply failed",
			give: notify.Entry{
				Message: "Compose apply failed. Docker Compose may have partially recreated services",
				Data: map[string]any{
					"container": "api",
					"error":     "compose apply exited with status 1",
				},
			},
			want: "Compose update failed for api. Docker Compose may have partially recreated services: " +
				"compose apply exited with status 1",
		},
		{
			name: "unreadable compose directory",
			give: notify.Entry{
				Message: "Compose project directory is not readable. Leaving the running container untouched",
				Data:    map[string]any{},
			},
			want: "Skipped container: compose directory is not readable",
		},
		{
			name: "invalid semver policy",
			give: notify.Entry{
				Message: "Compose apply skipped. Container has no compose service label",
				Data:    map[string]any{},
			},
			want: "Skipped container: container has no compose service label",
		},
		{
			name: "invalid git-host",
			give: notify.Entry{
				Message: "Skipped container with an invalid git-host",
				Data:    map[string]any{},
			},
			want: "Skipped container: invalid git-host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := notify.Data{Entries: []*notify.Entry{&tt.give}}

			result, err := Execute(templates.Templates["default-legacy"], payload, false)
			require.NoError(t, err)
			assert.Contains(t, result, tt.want)
			assert.NotContains(t, result, " | ")
		})
	}
}

func TestRenderDefaultLegacyEntriesRoot(t *testing.T) {
	t.Parallel()

	payload := notify.Data{
		Entries: []*notify.Entry{
			{
				Message: "Found new image",
				Data: map[string]any{
					"image":  "techwave/cyberscribe:latest",
					"new_id": "abc123def456",
				},
				Level: "info",
				Time:  time.Time{},
			},
		},
	}

	result, err := Execute(templates.Templates["default-legacy"], payload, false)
	require.NoError(t, err)
	assert.Contains(t, result, "Found new image: techwave/cyberscribe:latest (abc123def456)")
}

func TestRenderDefaultLegacyUsageCheckSkip(t *testing.T) {
	t.Parallel()

	payload := notify.Data{
		Entries: []*notify.Entry{
			{
				Message: "Failed to list containers for image usage check, skipping removal",
				Data: map[string]any{
					"error":      `Get "http://socket-proxy-write:2375/v1.55/containers/json?all=1": terminated signal received`,
					"image_id":   "769846626f2f",
					"image_name": "ghcr.io/amir20/dtop:latest",
				},
				Level: "warn",
				Time:  time.Time{},
			},
		},
	}

	result, err := Execute(templates.Templates["default-legacy"], payload, false)
	require.NoError(t, err)
	assert.Equal(
		t,
		"Skipped image cleanup: ghcr.io/amir20/dtop:latest (769846626f2f): "+
			`Get "http://socket-proxy-write:2375/v1.55/containers/json?all=1": terminated signal received`,
		result,
	)
	assert.NotContains(t, result, " | ")
}

func TestRenderDefaultLegacyDiskSpaceMessages(t *testing.T) {
	t.Parallel()

	payload := notify.Data{
		Entries: []*notify.Entry{
			{
				Message: "Docker image usage exceeds configured maximum",
				Data: map[string]any{
					"usage":       int64(32000000000),
					"max":         int64(40000000000),
					"warn":        int64(32000000000),
					"reclaimable": int64(4000000000),
					"image_count": int64(12),
				},
				Level: "error",
			},
			{
				Message: "Failed to query Docker image disk usage",
				Data:    map[string]any{"error": "daemon disk usage unavailable"},
				Level:   "error",
			},
		},
	}

	result, err := Execute(templates.Templates["default-legacy"], payload, false)
	require.NoError(t, err)
	assert.Contains(
		t,
		result,
		"Docker image usage exceeds configured maximum: 32 GB of 40 GB used (4 GB reclaimable, 12 images)",
	)
	assert.Contains(t, result, "Failed to query Docker image disk usage: daemon disk usage unavailable")
	assert.NotContains(t, result, " | ")
}

func TestRenderJSONIncludesReport(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["json.v1"],
		[]State{UpdatedState},
		[]LogLevel{InfoLevel},
	)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &decoded))
	require.Contains(t, decoded, "report")

	report, ok := decoded["report"].(map[string]any)
	require.True(t, ok)
	assert.NotEmpty(t, report["updated"])
	assert.Equal(t, "Title", decoded["title"])
	assert.Equal(t, "Host", decoded["host"])
}

func TestRenderLegacyModeOmitsReportRoot(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["json.v1"],
		nil,
		[]LogLevel{InfoLevel},
	)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(strings.TrimSpace(result), "["))
	assert.NotContains(t, result, `"report"`)
}

func TestRenderDeterministic(t *testing.T) {
	t.Parallel()

	input := templates.Templates["json.v1"]
	states := []State{UpdatedState}
	levels := []LogLevel{InfoLevel}

	first, err := Render(input, states, levels)
	require.NoError(t, err)

	second, err := Render(input, states, levels)
	require.NoError(t, err)

	assert.Equal(t, first, second)
}

func TestRenderDefaultTemplateWithoutReport(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["default"],
		nil,
		[]LogLevel{InfoLevel},
	)
	require.NoError(t, err)
	assert.NotContains(t, result, "can't evaluate field Report")
	assert.NotEmpty(t, result)
}

func TestRenderPreviewDefaultTemplateWithoutReport(t *testing.T) {
	t.Parallel()

	const previewDefault = `{{- if .Report -}}
 {{- with .Report -}}
   {{len .Scanned}} Scanned
 {{- end -}}
{{- else -}}
 {{range .Entries -}}{{.Message}}{{"\n"}}{{- end -}}
{{- end -}}`

	result, err := Render(previewDefault, nil, []LogLevel{InfoLevel})
	require.NoError(t, err)
	assert.NotContains(t, result, "can't evaluate field Report")
	assert.NotEmpty(t, result)
}

func TestRenderParseError(t *testing.T) {
	t.Parallel()

	_, err := Render("{{ .Report", nil, nil)
	require.Error(t, err)
}

func TestRenderPorcelainJSONWithNilReport(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["porcelain.json"],
		nil,
		nil,
	)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &decoded))
	require.Contains(t, decoded, "containers")

	containers, ok := decoded["containers"].([]any)
	require.True(t, ok)
	assert.Empty(t, containers)
}

func TestRenderPorcelainJSONWithContainers(t *testing.T) {
	t.Parallel()

	result, err := Render(
		templates.Templates["porcelain.json"],
		[]State{UpdatedState, FailedState, FreshState},
		nil,
	)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &decoded))

	containers, ok := decoded["containers"].([]any)
	require.True(t, ok)
	require.Len(t, containers, 3)

	first, ok := containers[0].(map[string]any)
	require.True(t, ok)
	assert.Contains(t, first, "name")
	assert.Contains(t, first, "image")
	assert.Contains(t, first, "image_id")
	assert.Contains(t, first, "latest_image_id")
	assert.Contains(t, first, "state")
	assert.Contains(t, first, "update_available")
}
