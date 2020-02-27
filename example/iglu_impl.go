package main

import (
	"net/http"
	"os"

	"github.com/Nacdlow/plugin-sdk"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// TestPlugin is an implementation of IgluPlugin
type TestPlugin struct {
	logger hclog.Logger
}

func (g *TestPlugin) OnLoad() error {
	g.logger.Debug("Loading test plugin!")
	return nil
}

func (g *TestPlugin) PluginHTTP(w http.ResponseWriter, req *http.Request) {
}

func (g *TestPlugin) GetManifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Id:      "test",
		Name:    "Test Plugin",
		Author:  "Nacdlow",
		Version: "v0.1.0",
	}
}

func (g *TestPlugin) RegisterDevice(reg sdk.DeviceRegistration) error {
	return nil
}

func (g *TestPlugin) OnDeviceToggle(id int, status bool) error {
	return nil
}

func (g *TestPlugin) GetWebExtensions() []sdk.WebExtension {
	return nil
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
		"iglu_plugin": &sdk.IgluPlugin{Impl: test},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
