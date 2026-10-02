package notify

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gowatchtower/tplprev/internal/report"
)

type stubReport struct {
	scanned []report.ContainerReport
	updated []report.ContainerReport
	failed  []report.ContainerReport
}

type stubContainerError struct {
	id    report.ContainerID
	name  string
	oldID report.ImageID
	newID report.ImageID
	image string
	state string
	err   string
	meta  stubGitMetadata
}

// stubGitMetadata carries the Git and OCI fields of a stub container.
//
// Fields are kept on the stub so tests can assert that populated metadata is
// emitted and empty metadata is omitted.
type stubGitMetadata struct {
	repo          string
	ref           string
	changelog     string
	source        string
	imageURL      string
	documentation string
	currentVer    string
	latestVer     string
	currentRev    string
	latestRev     string
}

func (r stubReport) Scanned() []report.ContainerReport   { return r.scanned }
func (r stubReport) Updated() []report.ContainerReport   { return r.updated }
func (r stubReport) Failed() []report.ContainerReport    { return r.failed }
func (r stubReport) Skipped() []report.ContainerReport   { return nil }
func (r stubReport) Stale() []report.ContainerReport     { return nil }
func (r stubReport) Fresh() []report.ContainerReport     { return nil }
func (r stubReport) Restarted() []report.ContainerReport { return nil }
func (r stubReport) All() []report.ContainerReport       { return r.scanned }

func (c stubContainerError) ID() report.ContainerID             { return c.id }
func (c stubContainerError) Name() string                       { return c.name }
func (c stubContainerError) CurrentImageID() report.ImageID     { return c.oldID }
func (c stubContainerError) LatestImageID() report.ImageID      { return c.newID }
func (c stubContainerError) ImageName() string                  { return c.image }
func (c stubContainerError) Error() string                      { return c.err }
func (c stubContainerError) State() string                      { return c.state }
func (c stubContainerError) IsMonitorOnly() bool                { return false }
func (c stubContainerError) NewContainerID() report.ContainerID { return "" }

func (c stubContainerError) GitRepo() string              { return c.meta.repo }
func (c stubContainerError) GitRef() string               { return c.meta.ref }
func (c stubContainerError) Changelog() string            { return c.meta.changelog }
func (c stubContainerError) Source() string               { return c.meta.source }
func (c stubContainerError) ImageURL() string             { return c.meta.imageURL }
func (c stubContainerError) Documentation() string        { return c.meta.documentation }
func (c stubContainerError) CurrentImageVersion() string  { return c.meta.currentVer }
func (c stubContainerError) LatestImageVersion() string   { return c.meta.latestVer }
func (c stubContainerError) CurrentImageRevision() string { return c.meta.currentRev }
func (c stubContainerError) LatestImageRevision() string  { return c.meta.latestRev }

func TestDataMarshalJSON(t *testing.T) {
	t.Parallel()

	updated := stubContainerError{
		id:    "sha256:c79110000000aaaa",
		name:  "updt1",
		oldID: "sha256:01d110000000aaaa",
		newID: "sha256:d0a110000000aaaa",
		image: "mock/updt1:latest",
		state: "updated",
	}
	failed := stubContainerError{
		id:    "sha256:c79210000000aaaa",
		name:  "fail1",
		oldID: "sha256:01d210000000aaaa",
		newID: "sha256:d0a210000000aaaa",
		image: "mock/fail1:latest",
		state: "failed",
		err:   "execution failed",
	}

	payload := Data{
		StaticData: StaticData{Title: "Title", Host: "Host"},
		Entries: []*Entry{
			{
				Message: "Found new image",
				Data:    map[string]any{"image": "mock/updt1:latest"},
				Time:    time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
				Level:   "info",
			},
		},
		Report: stubReport{
			scanned: []report.ContainerReport{updated, failed},
			updated: []report.ContainerReport{updated},
			failed:  []report.ContainerReport{failed},
		},
	}

	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(bytes, &decoded))

	assert.Equal(t, "Title", decoded["title"])
	assert.Equal(t, "Host", decoded["host"])

	report, ok := decoded["report"].(map[string]any)
	require.True(t, ok)
	assert.Len(t, report["scanned"], 2)
	assert.Len(t, report["updated"], 1)
	assert.Len(t, report["failed"], 1)

	failedJSON, ok := report["failed"].([]any)
	require.True(t, ok)
	failedEntry, ok := failedJSON[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "execution failed", failedEntry["error"])
	assert.Equal(t, "c79210000000", failedEntry["id"])

	entries, ok := decoded["entries"].([]any)
	require.True(t, ok)
	require.Len(t, entries, 1)
}

func TestDataMarshalJSONGitMetadata(t *testing.T) {
	t.Parallel()

	annotated := stubContainerError{
		id:    "sha256:c79310000000aaaa",
		name:  "git1",
		oldID: "sha256:01d310000000aaaa",
		newID: "sha256:d0a310000000aaaa",
		image: "mock/git1:latest",
		state: "updated",
		meta: stubGitMetadata{
			repo:          "https://github.com/org/app.git",
			ref:           "v1.4.0",
			changelog:     "https://github.com/org/app/releases/tag/v1.4.0",
			source:        "https://github.com/org/app",
			imageURL:      "https://github.com/org/app",
			documentation: "https://github.com/org/app#readme",
			currentVer:    "1.3.9",
			latestVer:     "1.4.0",
			currentRev:    "9f2c1ab4d7e5",
			latestRev:     "3d8b6e0f2c19",
		},
	}

	payload := Data{Report: stubReport{scanned: []report.ContainerReport{annotated}}}
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	var decoded struct {
		Report struct {
			Scanned []map[string]any `json:"scanned"`
		} `json:"report"`
	}
	require.NoError(t, json.Unmarshal(bytes, &decoded))
	require.Len(t, decoded.Report.Scanned, 1)

	entry := decoded.Report.Scanned[0]
	assert.Equal(t, "https://github.com/org/app.git", entry["gitRepo"])
	assert.Equal(t, "v1.4.0", entry["gitRef"])
	assert.Equal(t, "https://github.com/org/app/releases/tag/v1.4.0", entry["changelog"])
	assert.Equal(t, "https://github.com/org/app", entry["ociSource"])
	assert.Equal(t, "https://github.com/org/app", entry["imageUrl"])
	assert.Equal(t, "https://github.com/org/app#readme", entry["documentation"])
	assert.Equal(t, "1.3.9", entry["currentImageVersion"])
	assert.Equal(t, "1.4.0", entry["latestImageVersion"])
	assert.Equal(t, "9f2c1ab4d7e5", entry["currentImageRevision"])
	assert.Equal(t, "3d8b6e0f2c19", entry["latestImageRevision"])
}

func TestDataMarshalJSONOmitsEmptyGitMetadata(t *testing.T) {
	t.Parallel()

	plain := stubContainerError{
		id:    "sha256:c79410000000aaaa",
		name:  "plain1",
		oldID: "sha256:01d410000000aaaa",
		newID: "sha256:d0a410000000aaaa",
		image: "mock/plain1:latest",
		state: "updated",
	}

	payload := Data{Report: stubReport{scanned: []report.ContainerReport{plain}}}
	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	var decoded struct {
		Report struct {
			Scanned []map[string]any `json:"scanned"`
		} `json:"report"`
	}
	require.NoError(t, json.Unmarshal(bytes, &decoded))
	require.Len(t, decoded.Report.Scanned, 1)

	entry := decoded.Report.Scanned[0]
	for _, key := range []string{
		"gitRepo", "gitRef", "changelog", "ociSource", "imageUrl",
		"documentation", "currentImageVersion", "latestImageVersion",
		"currentImageRevision", "latestImageRevision",
	} {
		assert.NotContains(t, entry, key)
	}
}

func TestDataMarshalJSONNilReport(t *testing.T) {
	t.Parallel()

	payload := Data{
		StaticData: StaticData{Title: "Title", Host: "Host"},
		Entries:    []*Entry{},
	}

	bytes, err := json.Marshal(payload)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(bytes, &decoded))
	assert.Nil(t, decoded["report"])
}
