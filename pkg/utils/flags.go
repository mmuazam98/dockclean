package utils

import (
	"flag"
)

type Size struct {
	Value int64
	Unit  string
}

type Flags struct {
	DryRun        bool
	RemoveStopped bool
	VerboseMode   bool
	SizeLimit     Size
	B             bool
	KB            bool
	MB            bool
	GB            bool
}

func (f *Flags) GetSizeUnit() string {
	if f.B {
		return "b"
	} else if f.KB {
		return "kb"
	} else if f.MB {
		return "mb"
	} else if f.GB {
		return "gb"
	} else {
		return ""
	}
}

func ParseFlags() *Flags {
	f := &Flags{}
	flag.BoolVar(&f.DryRun, "dry-run", false, "List unused Docker images without deleting them")
	flag.BoolVar(&f.RemoveStopped, "remove-stopped", false, "Remove Images Associated with Stopped Containers")
	flag.BoolVar(&f.VerboseMode, "verbose", false, "Verbose mode provides additional details about each image during cleanup")
	flag.Int64Var(&f.SizeLimit.Value, "size-limit", 0, "Specify the size limit to filter images (e.g., 500MB, 1GB)")
	flag.BoolVar(&f.B, "B", false, "Specify the size unit as bytes")
	flag.BoolVar(&f.KB, "KB", false, "Specify the size unit as kilobytes")
	flag.BoolVar(&f.MB, "MB", false, "Specify the size unit as megabytes")
	flag.BoolVar(&f.GB, "GB", false, "Specify the size unit as gigabytes")

	flag.Parse()
	return f
}
