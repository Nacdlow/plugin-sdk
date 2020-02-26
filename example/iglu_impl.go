package main

import (
	"os"

	"github.com/Nacdlow/plugin-sdk"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// TestPlugin is an implementation of IgluPlugin
type TestPlugin struct {
	logger hclog.Logger
}

func (g *TestPlugin) OnLoad() {
	g.logger.Debug("Loading test plugin!")
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "IGLU_PLUGIN",
	MagicCookieValue: "MzlK0OGpIRs",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	test := &TestPlugin{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"iglu_plugin": &api.IgluPlugin{Impl: test},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}