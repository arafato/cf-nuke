package version

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"text/template"
)

// Injected at build time via -ldflags
var (
	version = ""
	commit  = ""
	date    = ""
)

const versionTemplateLong = `cf-nuke {{.Version}}
Commit:     {{.Commit}}
Built:      {{.Date}}
Go version: {{.GoVersion}}
OS/Arch:    {{.Os}}/{{.Arch}}
`

// Print outputs the full version information
func Print(wr io.Writer) error {
	v, c, d := getVersionInfo()

	tmpl, err := template.New("version").Parse(versionTemplateLong)
	if err != nil {
		return err
	}

	data := struct {
		Version   string
		Commit    string
		Date      string
		GoVersion string
		Os        string
		Arch      string
	}{
		Version:   v,
		Commit:    c,
		Date:      d,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	return tmpl.Execute(wr, data)
}

// PrintShort outputs a single-line version string
func PrintShort(wr io.Writer) error {
	v, c, d := getVersionInfo()
	// Extract just the date part (YYYY-MM-DD) from ISO format
	dateShort := d
	if len(d) >= 10 {
		dateShort = d[:10]
	}
	_, err := fmt.Fprintf(wr, "cf-nuke %s (%s, %s)\n", v, c, dateShort)
	return err
}

// GetVersion returns just the version string
func GetVersion() string {
	v, _, _ := getVersionInfo()
	return v
}

func getVersionInfo() (string, string, string) {
	// If injected at build time, use those values
	if version != "" && commit != "" && date != "" {
		return version, commit, date
	}

	// Fallback: extract from Go build info (works with `go install`)
	v, c, d := "dev", "unknown", "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if len(s.Value) >= 7 {
					c = s.Value[:7]
				} else {
					c = s.Value
				}
			case "vcs.time":
				d = s.Value
			}
		}
	}
	return v, c, d
}
