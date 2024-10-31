package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
)

// Helper function to truncate docker image ID with ellipsis
func FormatDockerImageID(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper function to converts bytes to human-readable format
func FormatSize(bytes int64) string {
	const unit = 1024.0
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper function used to format the docker image labels Key Value Pairs
func FormatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "No Labels Found"
	}

	var labelStore []string
	for labelKey, labelValue := range labels {
		// Handling empty values
		if labelValue == "" {
			labelStore = append(labelStore, labelKey)
			continue
		}
		labelStore = append(labelStore, fmt.Sprintf("%s:%s", labelKey, labelValue))
	}

	return strings.Join(labelStore, ", ")
}

// converts value to bytes
func ToBytes(size float64, unit string) float64 {
	multipliers := map[string]float64{
		"B":  1.0,
		"KB": 1024.0,
		"MB": 1024.0 * 1024.0,
		"GB": 1024.0 * 1024.0 * 1024.0,
	}

	return size * multipliers[unit]
}

// RemoveDockerImage
func RemoveDockerImage(d *DockerClient, ctx context.Context, imageID string, opts image.RemoveOptions) error {
	_, err := d.CLI.ImageRemove(context.Background(), imageID, opts)
	return err
}

// Measure Execution Time
func MeasureExecutionTime(f func()) string {
	startTime := time.Now()
	f()
	duration := time.Since(startTime)
	return duration.String()
}
