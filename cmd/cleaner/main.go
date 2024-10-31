package main

import (
	"fmt"
	"log"

	"github.com/mmuazam98/dockclean/internal/docker"
	"github.com/mmuazam98/dockclean/pkg/utils"
)

func main() {

	var dc docker.DockerClient

	err := dc.InitDockerClient()
	if err != nil {
		log.Fatalf("Error initializing docker client : %v", err)
	}

	f := utils.ParseFlags()

	switch {
	case f.DryRun:
		dc.PrintUnusedImages()
	case f.RemoveStopped:
		dc.CleanupStoppedContainerImages()
	case f.VerboseMode:
		dc.VerboseModeCleanup()
	case f.SizeLimit.Value >= 0:
		var unit string = f.GetSizeUnit()
		if unit == "" {
			log.Fatalf("Please specify a size unit (B, KB, MB, or GB)")
		}
		dc.RemoveExceedSizeLimit(f.SizeLimit.Value, unit)
	default:
		var exectionTime string
		if f.ConcurrentDelete {
			exectionTime = docker.MeasureExecutionTime(func() { dc.RemoveUnusedImages(true) })
			fmt.Println("Execution Time ( Concurrent ) : ", exectionTime)
		} else {
			exectionTime = docker.MeasureExecutionTime(func() { dc.RemoveUnusedImages(false) })
			fmt.Println("Execution Time ( Sequential ) : ", exectionTime)
		}

	}

}
