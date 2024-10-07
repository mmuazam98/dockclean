package utils

import "flag"

type Flags struct {
	DryRun        bool
	RemoveStopped bool
}

func ParseFlags() *Flags {
	f := &Flags{}
	flag.BoolVar(&f.DryRun, "dry-run", false, "List unused Docker images without deleting them")
	flag.BoolVar(&f.RemoveStopped, "remove-stopped", false, "Remove Images Associated with Stopped Containers")
	flag.Parse()
	return f
}
