package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v68/github"
	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v3"
)

var logLevel = slog.LevelInfo

func main() {
	programLevel := new(slog.LevelVar)
	programLevel.Set(logLevel)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})))

	ctx := context.Background()

	// Step 1: Get the latest release tag
	tag, err := getCFTag(ctx)
	if err != nil {
		slog.Error("Failed to fetch the latest release tag", "error", err)
		os.Exit(1)
	}
	slog.Info("Fetched the latest release tag", "tag", tag)

	// Step 2: Check if the tag exists in Docker Registry
	client := &http.Client{}
	exists, err := checkTagExists(ctx, tag, client)
	if err != nil {
		slog.Error("Error checking tag existence in Docker Registry", "error", err)
		os.Exit(1)
	}

	if !exists {
		slog.Warn("Tag does not exist in Docker Registry", "tag", tag)
		os.Exit(0)
	}
	slog.Info("Tag exists in Docker Registry", "tag", tag)

	// Step 3: Read and update Chart.yaml
	chartPath := os.Getenv("CHART_PATH")
	if chartPath == "" {
		chartPath = "charts/cloudflare-tunnel/Chart.yaml" // Default path
	}

	if err := updateChartYAML(chartPath, tag); err != nil {
		slog.Error("Error updating Chart.yaml", "error", err)
		os.Exit(1)
	}
	slog.Info("Updated Chart.yaml successfully", "newAppVersion", tag)
}

// Fetch the latest GitHub release tag
func getCFTag(ctx context.Context) (string, error) {
	ghClient := github.NewClient(nil)
	owner := "cloudflare"
	repo := "cloudflared"

	release, _, err := ghClient.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch latest release")
	}

	return release.GetTagName(), nil
}

// Check if a Docker tag exists
func checkTagExists(ctx context.Context, tag string, client *http.Client) (bool, error) {
	url := "https://registry.hub.docker.com/v2/repositories/cloudflare/cloudflared/tags/" + tag
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, errors.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
}

// Update Chart.yaml with the new appVersion and bumped version
func updateChartYAML(path, newAppVersion string) error {
	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open Chart.yaml")
	}
	defer file.Close()

	// Parse the YAML file into a map
	var data map[string]interface{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return errors.Wrap(err, "failed to parse Chart.yaml")
	}

	// Update the appVersion
	data["appVersion"] = newAppVersion

	// Update the version (bump minor)
	currentVersion, ok := data["version"].(string)
	if !ok {
		return errors.New("version field not found or invalid in Chart.yaml")
	}

	newVersion, err := bumpMinorVersion(currentVersion)
	if err != nil {
		return errors.Wrap(err, "failed to bump minor version")
	}
	data["version"] = newVersion

	// Update annotations
	annotations, ok := data["annotations"].(map[string]interface{})
	if !ok {
		annotations = make(map[string]interface{})
	}
	changeLog := fmt.Sprintf("- kind: changed\n  description: Bump app version to %s\n", newAppVersion)
	annotations["artifacthub.io/changes"] = changeLog
	data["annotations"] = annotations

	// Write the updated YAML back to the file
	outputFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create Chart.yaml for writing")
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	defer encoder.Close()
	if err := encoder.Encode(data); err != nil {
		return errors.Wrap(err, "failed to write updated Chart.yaml")
	}

	return nil
}

// Bump the minor version of a semantic version string
func bumpMinorVersion(currentVersion string) (string, error) {
	ver, err := version.NewVersion(currentVersion)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse version")
	}

	segments := ver.Segments()
	if len(segments) < 2 {
		return "", errors.New("version must have at least major and minor segments")
	}

	segments[1]++ // Increment minor version
	if len(segments) > 2 {
		segments[2] = 0 // Reset patch version
	}

	return fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2]), nil
}
