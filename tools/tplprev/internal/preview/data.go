package preview

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gowatchtower/tplprev/internal/notify"
	"github.com/gowatchtower/tplprev/internal/report"
)

// PreviewData represents a generator for preview data, including container statuses and log entries.
type PreviewData struct {
	rand           *rand.Rand
	lastTime       time.Time
	report         *sessionReport
	containerCount int
	entries        []*notify.Entry
	static         notify.StaticData
}

const (
	idLength                = 32
	maxTimeIncrementSeconds = 30
	previewInstanceCount    = 2
	previewRemovedInstances = 1
	previewAPIPort          = "8080"
	previewCooldown         = "24h"
	previewEligibleIn       = "6h"
	previewEligibleAfter    = 6 * time.Hour
	previewDiskUsage        = int64(32_000_000_000)
	previewDiskMax          = int64(40_000_000_000)
	previewDiskWarn         = int64(32_000_000_000)
	previewDiskReclaimable  = int64(4_000_000_000)
	previewDiskImageCount   = int64(12)
	previewGitRepo          = "https://github.com/gowatchtower/app.git"
	previewGitRef           = "v1.4.0"
	previewChangelog        = "https://github.com/gowatchtower/app/releases/tag/v1.4.0"
	previewCurrentVersion   = "1.3.9"
	previewLatestVersion    = "1.4.0"
	previewRevision         = "9f2c1ab4d7e5"
	previewNextRevision     = "3d8b6e0f2c19"
	previewImageURL         = "https://github.com/gowatchtower/app"
	previewDocumentationURL = "https://github.com/gowatchtower/app#readme"
)

// gitMetadata describes the Git and OCI metadata attached to a preview container.
//
// A preview leaves every field empty most of the time so the default templates
// render the same output as an ordinary registry update. Metadata is only
// populated for containers the generator marks as Git-associated or OCI
// annotated, which mirrors how Watchtower resolves these fields.
type gitMetadata struct {
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

var (
	errExecutionFailed = errors.New("execution failed")
	errSkipped         = errors.New("container skipped")

	// previewStartTime is the fixed session start so generated timestamps stay stable.
	previewStartTime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
)

// New initializes a PreviewData generator with a fixed seed for deterministic output.
func New() *PreviewData {
	//nolint:gosec // math/rand is used for deterministic preview output.
	return &PreviewData{
		rand:           rand.New(rand.NewSource(1)),
		lastTime:       previewStartTime,
		report:         &sessionReport{},
		containerCount: 0,
		entries:        []*notify.Entry{},
		static: notify.StaticData{
			Title: "Title",
			Host:  "Host",
		},
	}
}

// AddFromState adds a container status entry to the report with the specified state.
//
// Parameters:
//   - state: Container report state to simulate.
func (p *PreviewData) AddFromState(state State) {
	cid := report.ContainerID(p.generateID())
	oldImageID := report.ImageID(p.generateID())
	newImageID := report.ImageID(p.generateID())
	name := p.generateName()
	image := p.generateImageName(name)

	var err error

	//nolint:exhaustive // Only failed and skipped states carry a simulated error.
	switch state {
	case FailedState:
		err = fmt.Errorf("%w: %s", errExecutionFailed, p.randomEntry(errorMessages))
	case SkippedState:
		err = fmt.Errorf("%w: %s", errSkipped, p.randomEntry(skippedMessages))
	}

	status := &containerStatus{
		containerID:    cid,
		oldImage:       oldImageID,
		newImage:       newImageID,
		containerName:  name,
		imageName:      image,
		containerError: err,
		state:          state,
	}
	applyGitMetadata(status, state, p.containerCount)

	switch state {
	case ScannedState:
		p.report.scanned = append(p.report.scanned, status)
	case UpdatedState:
		p.report.updated = append(p.report.updated, status)
	case FailedState:
		p.report.failed = append(p.report.failed, status)
	case SkippedState:
		p.report.skipped = append(p.report.skipped, status)
	case RestartedState:
		p.report.restarted = append(p.report.restarted, status)
	case StaleState:
		p.report.stale = append(p.report.stale, status)
	case FreshState:
		p.report.fresh = append(p.report.fresh, status)
	}

	p.containerCount++
}

// applyGitMetadata populates the Git and OCI report fields on a preview container.
//
// Updated containers carry the full set so the default report template shows a
// ref and a changelog. Fresh containers carry only the fields describing the
// image already running, matching a registry update that found nothing new.
// Every other state leaves the fields empty, as Watchtower resolves them only
// for containers it actually reports on.
//
// Parameters:
//   - status: Container report to populate.
//   - state: State the container was added with.
//   - index: Container index, used to vary the preview across several containers.
func applyGitMetadata(status *containerStatus, state State, index int) {
	var meta gitMetadata

	//nolint:exhaustive // Only updated and fresh containers resolve Git and OCI metadata.
	switch state {
	case UpdatedState:
		meta = gitMetadata{
			repo:          previewGitRepo,
			ref:           previewGitRef,
			changelog:     previewChangelog,
			source:        previewImageURL,
			imageURL:      previewImageURL,
			documentation: previewDocumentationURL,
			currentVer:    previewCurrentVersion,
			latestVer:     previewLatestVersion,
			currentRev:    previewRevision,
			latestRev:     previewNextRevision,
		}
	case FreshState:
		meta = gitMetadata{
			repo:          previewGitRepo,
			source:        previewImageURL,
			imageURL:      previewImageURL,
			documentation: previewDocumentationURL,
			currentVer:    previewCurrentVersion,
			currentRev:    previewRevision,
		}
	default:
		return
	}

	// Vary the second and later containers so a multi-container preview does not
	// repeat identical metadata on every line.
	if index > 0 {
		suffix := strconv.FormatInt(int64(index), 10)
		meta.repo = meta.repo + "/" + suffix
		meta.changelog = meta.changelog + "-" + suffix
	}

	status.gitRepo = meta.repo
	status.gitRef = meta.ref
	status.changelog = meta.changelog
	status.source = meta.source
	status.imageURL = meta.imageURL
	status.documentation = meta.documentation
	status.currentVersion = meta.currentVer
	status.latestVersion = meta.latestVer
	status.currentRev = meta.currentRev
	status.latestRev = meta.latestRev
}

// AddLogEntry adds a preview log entry with the specified level.
//
// Parameters:
//   - level: Log level whose message pool should be sampled.
func (p *PreviewData) AddLogEntry(level LogLevel) {
	var msg string

	switch level {
	case FatalLevel, ErrorLevel:
		msg = p.randomEntry(errorMessages)
	case WarnLevel:
		msg = p.randomEntry(warningMessages)
	case TraceLevel, DebugLevel, InfoLevel, PanicLevel:
		msg = p.randomEntry(infoMessages)
	default:
		msg = p.randomEntry(infoMessages)
	}

	name, image := p.logSubject()

	p.entries = append(p.entries, &notify.Entry{
		Message: msg,
		Data:    p.dataForMessage(msg, name, image),
		Time:    p.generateTime(),
		Level:   level.String(),
	})
}

// NotificationData returns the generated notification template data.
//
// Report is nil when no container states were added so templates treat the
// preview as legacy/log mode.
//
// Returns:
//   - notify.Data: Simulated notification data.
func (p *PreviewData) NotificationData() notify.Data {
	payload := notify.Data{
		StaticData: p.static,
		Entries:    p.entries,
	}

	if p.containerCount > 0 {
		payload.Report = p.report
	}

	return payload
}

func (p *PreviewData) generateID() string {
	buf := make([]byte, idLength)
	_, _ = p.rand.Read(buf)

	return hex.EncodeToString(buf)
}

func (p *PreviewData) generateTime() time.Time {
	p.lastTime = p.lastTime.Add(time.Duration(p.rand.Intn(maxTimeIncrementSeconds)) * time.Second)

	return p.lastTime
}

func (p *PreviewData) randomEntry(arr []string) string {
	return arr[p.rand.Intn(len(arr))]
}

func (p *PreviewData) generateName() string {
	index := p.containerCount
	if index < len(containerNames) {
		return containerNames[index]
	}

	suffix := index / len(containerNames)
	index %= len(containerNames)

	return containerNames[index] + strconv.FormatInt(int64(suffix), 10)
}

func (p *PreviewData) generateImageName(name string) string {
	index := p.containerCount % len(organizationNames)

	return organizationNames[index] + "/" + name + ":latest"
}

func (p *PreviewData) logSubject() (string, string) {
	index := len(p.entries)
	name := containerNames[index%len(containerNames)]
	org := organizationNames[index%len(organizationNames)]

	return name, org + "/" + name + ":latest"
}

func (p *PreviewData) dataForMessage(msg, name, image string) map[string]any {
	cid := report.ContainerID(p.generateID()).ShortID()
	newID := report.ImageID(p.generateID()).ShortID()
	oldID := report.ImageID(p.generateID()).ShortID()

	switch msg {
	case "Found new image":
		return map[string]any{"image": image, "new_id": newID}
	case "Stopping container", "Stopping linked container":
		return map[string]any{"container": name, "id": cid}
	case "Started new container", "Started linked container":
		return map[string]any{"container": name, "new_id": newID}
	case "Removing image":
		return map[string]any{"image_name": image, "image_id": oldID}
	case "Failed to list containers for image usage check, skipping removal":
		return map[string]any{
			"image_name": image,
			"image_id":   oldID,
			"error":      "cannot verify image usage",
		}
	case "Container updated":
		return map[string]any{
			"container": name,
			"image":     image,
			"old_id":    oldID,
			"new_id":    newID,
		}
	case "Detected multiple Watchtower instances - initiating cleanup":
		return map[string]any{"count": previewInstanceCount}
	case "Successfully removed all excess Watchtower containers":
		return map[string]any{"removed_instances": previewRemovedInstances}
	case "Image is within cooldown period - not eligible for update":
		return map[string]any{
			"image":       image,
			"cooldown":    previewCooldown,
			"eligible_in": previewEligibleIn,
			"eligible_at": p.lastTime.Add(previewEligibleAfter).Format(time.RFC3339),
		}
	case "Image age exceeds cooldown - eligible for update",
		"Image creation time unavailable - update check unavailable":
		return map[string]any{"image": image, "cooldown": previewCooldown}
	case "Starting HTTP API server", "HTTP API server is enabled":
		return map[string]any{"tls": false, "host": "0.0.0.0", "port": previewAPIPort}
	case "Only checking containers in scope":
		return map[string]any{"scope": "prod"}
	case "Docker image usage exceeds configured maximum":
		return map[string]any{
			"usage":       previewDiskMax,
			"max":         previewDiskMax,
			"warn":        previewDiskWarn,
			"reclaimable": previewDiskReclaimable,
			"image_count": previewDiskImageCount,
		}
	case "Docker image usage exceeds configured warning threshold":
		return map[string]any{
			"usage":       previewDiskUsage,
			"max":         previewDiskMax,
			"warn":        previewDiskWarn,
			"reclaimable": previewDiskReclaimable,
			"image_count": previewDiskImageCount,
		}
	case "Failed to query Docker image disk usage":
		return map[string]any{"error": "daemon disk usage unavailable"}
	case "Changelog":
		return map[string]any{"changelog": previewChangelog}
	case "Found new Git revision":
		return map[string]any{
			"repo":         previewGitRepo,
			"ref":          previewGitRef,
			"revision":     previewRevision,
			"short_commit": previewRevision[:12],
			"changelog":    previewChangelog,
		}
	case "Built image from Git URL context":
		return map[string]any{
			"image":     organizationNames[p.containerCount%len(organizationNames)] + "/" + name + ":git-" + previewRevision[:12],
			"image_id":  newID,
			"repo":      previewGitRepo,
			"commit":    previewRevision,
			"changelog": previewChangelog,
		}
	case "Built Compose project":
		return map[string]any{"project": "webstack", "service": name, "dir": "/srv/webstack"}
	case "Git build failed. Leaving running container untouched":
		return map[string]any{
			"container": name,
			"image":     image,
			"repo":      previewGitRepo,
			"commit":    previewRevision,
			"error":     "docker build failed",
		}
	case "Compose apply failed. Docker Compose may have partially recreated services":
		return map[string]any{"container": name, "project": "webstack", "error": "compose apply exited with status 1"}
	case "Compose project directory is not readable. Leaving the running container untouched":
		return map[string]any{"container": name, "dir": "/srv/webstack"}
	case "Skipped container with an invalid git semver policy":
		return map[string]any{"container": name, "repo": previewGitRepo}
	case "Skipped container with an invalid git-host":
		return map[string]any{"container": name, "repo": previewGitRepo}
	case "Docker image usage budget enabled":
		return map[string]any{
			"disk_space_max":  previewDiskMax,
			"disk_space_warn": previewDiskWarn,
		}
	default:
		return nil
	}
}
