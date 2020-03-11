package sdk

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// PluginManifest is used to describe the plugin's id, name, author, version, etc.
type PluginManifest struct {
	Id, Name, Author, Version string
}

// ExtensionType specifies the type of web extension.
type ExtensionType int

const (
	CSS = iota + 1
	JavaScript
)

// WebExtension represents an addon to the web page.
type WebExtension struct {
	Type           ExtensionType
	PathMatchRegex string
	Source         string
}

const (
	StringValue = iota
	NumberValue
	OptionValue // drop down list
	BooleanValue
)

// PluginConfig represents a plugin's configuration key-value item constraints
// including title and description.
type PluginConfig struct {
	Title          string
	Description    string
	Key            string // a unique key
	Type           int
	Default        string // Not for OptionValues
	Placeholder    string
	Values         []string // for option lists, first value is default
	IsUserSpecific bool
}

// ConfigKV represents a set key-value configuration, used in communicating to
// the plugin the current configuration.
type ConfigKV struct {
	Key    string
	Value  string
	UserID string
}

// AvailableDevice represents an available device the plugin may be able to
// register.
type AvailableDevice struct {
	UniqueID         string
	ManufacturerName string
	ModelName        string
	Type             int64
}

// Iglu is the interface that we're exposing as a plugin.
type Iglu interface {
	OnLoad() error
	GetManifest() PluginManifest
	OnDeviceToggle(uniqueID string, status bool)
	GetDeviceStatus(uniqueID string) bool
	GetWebExtensions() []WebExtension
	GetPluginConfiguration() []PluginConfig
	OnConfigurationUpdate(conf []ConfigKV)
	GetAvailableDevices() []AvailableDevice
}

// IgluRPC is what the server is using to communicate to the plugin over RPC
type IgluRPC struct {
	client *rpc.Client
}

type OnLoadReply struct {
	Err error
}

func (i *IgluRPC) OnLoad() error {
	rep := OnLoadReply{}
	err := i.client.Call("Plugin.OnLoad", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Err
}

type GetManifestReply struct {
	Manifest PluginManifest
}

func (i *IgluRPC) GetManifest() PluginManifest {
	rep := GetManifestReply{}
	err := i.client.Call("Plugin.GetManifest", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Manifest
}

type OnDeviceToggleArgs struct {
	Id     string
	Status bool
}

func (i *IgluRPC) OnDeviceToggle(id string, status bool) {
	args := &OnDeviceToggleArgs{Id: id, Status: status}
	err := i.client.Call("Plugin.OnDeviceToggle", args, 0)
	if err != nil {
		panic(err)
	}
}

type GetWebExtensionsReply struct {
	Extensions []WebExtension
}

func (i *IgluRPC) GetWebExtensions() []WebExtension {
	rep := &GetWebExtensionsReply{}
	err := i.client.Call("Plugin.GetWebExtensions", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Extensions
}

type GetPluginConfigurationReply struct {
	Configuration []PluginConfig
}

func (i *IgluRPC) GetPluginConfiguration() []PluginConfig {
	rep := &GetPluginConfigurationReply{}
	err := i.client.Call("Plugin.GetPluginConfiguration", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Configuration
}

type OnConfigurationUpdateArgs struct {
	Configuration []ConfigKV
}

func (i *IgluRPC) OnConfigurationUpdate(config []ConfigKV) {
	args := &OnConfigurationUpdateArgs{Configuration: config}
	err := i.client.Call("Plugin.GetPluginConfiguration", args, 0)
	if err != nil {
		panic(err)
	}
	return
}

type GetAvailableDevicesReply struct {
	Devices []AvailableDevice
}

func (i *IgluRPC) GetAvailableDevices() []AvailableDevice {
	rep := &GetAvailableDevicesReply{}
	err := i.client.Call("Plugin.GetAvailableDevices", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Devices
}

type GetDeviceStatusArgs struct {
	UniqueID string
}

type GetDeviceStatusReply struct {
	Status bool
}

func (i *IgluRPC) GetDeviceStatus(uniqueID string) bool {
	args := &GetDeviceStatusArgs{UniqueID: uniqueID}
	rep := &GetDeviceStatusReply{}
	err := i.client.Call("Plugin.GetDeviceStatus", args, &rep)
	if err != nil {
		panic(err)
	}
	return rep.Status
}

// This is the implementation of plugin.Plugin.
type IgluPlugin struct {
	Impl Iglu
}

func (p *IgluPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &IgluRPCServer{Impl: p.Impl}, nil
}

func (IgluPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &IgluRPC{client: c}, nil
}
