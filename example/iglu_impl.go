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

func (g *TestPlugin) OnLoad() error {
	g.logger.Debug("Loading test plugin!")
	return nil
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

func (g *TestPlugin) GetPluginConfiguration() []sdk.PluginConfig {
	return []sdk.PluginConfig{
		{Title: "Option Text", Description: "This is a text field.", Key: "text1",
			Type: sdk.StringValue, IsUserSpecific: true},
		{Title: "Option Bool", Description: "This is a true/false field.", Key: "bool1",
			Type: sdk.BooleanValue, IsUserSpecific: false},
		{Title: "Option Num", Description: "This is a number field.", Key: "num1",
			Type: sdk.NumberValue, IsUserSpecific: false},
		{Title: "Option Selection", Description: "This is a selection field..", Key: "select1",
			Type: sdk.OptionValue, Values: []string{"Default", "Easy", "Hard"}, IsUserSpecific: true},
	}
}

func (g *TestPlugin) OnConfigurationUpdate(config []sdk.ConfigKV) {
}

func (g *TestPlugin) GetAvailableDevices() []sdk.AvailableDevice {
	return []sdk.AvailableDevice{
		{UniqueID: "TestDevice01", ManufacturerName: "Acme Corp", ModelName: "Acme LED Hue", Type: 0},
	}
}

func (g *TestPlugin) GetWebExtensions() []sdk.WebExtension {
	return []sdk.WebExtension{
		{Type: sdk.CSS, Source: "/*Test CSS*/", PathMatchRegex: "*"},
		{Type: sdk.JavaScript, Source: "/*Test JS*/", PathMatchRegex: "*"},
	}
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
