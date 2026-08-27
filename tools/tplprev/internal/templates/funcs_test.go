package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gowatchtower/tplprev/internal/report"
)

func TestHasKey(t *testing.T) {
	t.Parallel()

	data := map[string]any{"usage": int64(0)}
	assert.True(t, hasKey(data, "usage"))
	assert.False(t, hasKey(data, "max"))
	assert.False(t, hasKey("not-a-map", "usage"))
	assert.False(t, hasKey(nil, "usage"))
}

func TestToJSON(t *testing.T) {
	t.Parallel()

	got := toJSON(map[string]any{"title": "Title"})
	assert.Contains(t, got, `"title": "Title"`)
}

func TestToPorcelainJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give any
		want string
	}{
		{
			name: "nil report",
			give: nil,
			want: `{
  "containers": []
}`,
		},
		{
			name: "wrong type",
			give: "not a report",
			want: "failed to marshal porcelain JSON: input is not a report.Report",
		},
		{
			name: "empty report",
			give: stubReport{},
			want: `{
  "containers": []
}`,
		},
		{
			name: "updated and failed containers",
			give: stubReport{
				containers: []stubContainer{
					{
						name:         "nginx",
						image:        "nginx:latest",
						currentID:    "abc123def4567890",
						latestID:     "def456abc1237890",
						state:        "Updated",
						containerErr: "",
					},
					{
						name:         "fail1",
						image:        "mock/fail1:latest",
						currentID:    "01d210000000abcd",
						latestID:     "d0a210000000abcd",
						state:        "Failed",
						containerErr: "accidentally the whole container",
					},
					{
						name:         "frsh1",
						image:        "mock/frsh1:latest",
						currentID:    "01d310000000abcd",
						latestID:     "01d310000000abcd",
						state:        "Fresh",
						containerErr: "",
					},
				},
			},
			want: `{
  "containers": [
    {
      "name": "nginx",
      "image": "nginx:latest",
      "image_id": "abc123def456",
      "latest_image_id": "def456abc123",
      "state": "Updated",
      "update_available": true
    },
    {
      "name": "fail1",
      "image": "mock/fail1:latest",
      "image_id": "01d210000000",
      "latest_image_id": "d0a210000000",
      "state": "Failed",
      "update_available": true,
      "error": "accidentally the whole container"
    },
    {
      "name": "frsh1",
      "image": "mock/frsh1:latest",
      "image_id": "01d310000000",
      "latest_image_id": "01d310000000",
      "state": "Fresh",
      "update_available": false
    }
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, toPorcelainJSON(tt.give))
		})
	}
}

type stubReport struct {
	containers []stubContainer
}

func (s stubReport) Scanned() []report.ContainerReport   { return nil }
func (s stubReport) Updated() []report.ContainerReport   { return nil }
func (s stubReport) Failed() []report.ContainerReport    { return nil }
func (s stubReport) Skipped() []report.ContainerReport   { return nil }
func (s stubReport) Stale() []report.ContainerReport     { return nil }
func (s stubReport) Fresh() []report.ContainerReport     { return nil }
func (s stubReport) Restarted() []report.ContainerReport { return nil }

func (s stubReport) All() []report.ContainerReport {
	all := make([]report.ContainerReport, 0, len(s.containers))
	for i := range s.containers {
		all = append(all, s.containers[i])
	}

	return all
}

//nolint:errname // stubContainer is not an error type. It implements ContainerReport.Error.
type stubContainer struct {
	name         string
	image        string
	currentID    report.ImageID
	latestID     report.ImageID
	state        string
	containerErr string
}

func (s stubContainer) ID() report.ContainerID             { return "" }
func (s stubContainer) Name() string                       { return s.name }
func (s stubContainer) CurrentImageID() report.ImageID     { return s.currentID }
func (s stubContainer) LatestImageID() report.ImageID      { return s.latestID }
func (s stubContainer) ImageName() string                  { return s.image }
func (s stubContainer) Error() string                      { return s.containerErr }
func (s stubContainer) State() string                      { return s.state }
func (s stubContainer) IsMonitorOnly() bool                { return false }
func (s stubContainer) NewContainerID() report.ContainerID { return "" }

func TestFormatRFC1123(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give string
		want string
	}{
		{give: "2026-01-02T03:04:05Z", want: "Fri, 02 Jan 2026 03:04:05 UTC"},
		{give: "not-a-timestamp", want: "not-a-timestamp"},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, formatRFC1123(tt.give))
		})
	}
}
