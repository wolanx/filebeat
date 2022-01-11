package main

import (
	"os"

	"github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
	_ "github.com/wolanx/filebeat/pkg/loki"
)

// go run main.go -e -c filebeat.yml
func main() {
	if err := cmd.Filebeat(inputs.Init, cmd.FilebeatSettings()).Execute(); err != nil {
		os.Exit(1)
	}
}
