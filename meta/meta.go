package meta

import "os"

var (
	Version     string
	GoVersion   string
	GitCommit   string
	Built       string
	OsArch      string
	ENV         string
	Hostname, _ = os.Hostname()
)
