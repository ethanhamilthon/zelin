package version

import "runtime/debug"

const Name = "Zelin"

var (
	Number = "0.1.0"
	Commit = ""
	Date   = ""
)

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func Get() Info {
	info := Info{Name: Name, Version: Number, Commit: Commit, Date: Date}
	if build, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range build.Settings {
			if setting.Key == "vcs.revision" && info.Commit == "" {
				info.Commit = setting.Value
			}
			if setting.Key == "vcs.time" && info.Date == "" {
				info.Date = setting.Value
			}
		}
	}
	if len(info.Commit) > 12 {
		info.Commit = info.Commit[:12]
	}
	return info
}
