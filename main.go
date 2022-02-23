package main

import (
	"runtime"
	"strings"
	"time"

	server "grafana-adapter/modules/router"
	"grafana-adapter/modules/settings"
)

var (
	// Version holds the current grafana-adapter version
	Version = "development"
	// Tags holds the build tags used
	Tags = ""
	// MakeVersion holds the current Make version if built with make
	MakeVersion = ""
)

func formatBuiltWith() string {
	var version = runtime.Version()
	if len(MakeVersion) > 0 {
		version = MakeVersion + ", " + runtime.Version()
	}
	if len(Tags) == 0 {
		return " built with " + version
	}

	return " built with " + version + " : " + strings.ReplaceAll(Tags, " ", ", ")
}

func init() {
	settings.AppVer = Version
	settings.AppBuiltWith = formatBuiltWith()
	settings.AppStartTime = time.Now().UTC()
}

func main() {
	server.Start()
}
