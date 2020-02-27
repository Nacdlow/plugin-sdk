package sdk

import (
	"net/http"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// PluginManifest is used to describe the plugin's id, name, author, version, etc.
type PluginManifest struct {
	Id, Name, Author, Version string
}

// DeviceRegistration represents a device to be registered to a plugin. This is
// used to inform a plugin about a device.
type DeviceRegistration struct {
	DeviceID    int
	Description string
	Type        int64 // TODO use enum
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

// Iglu is the interface that we're exposing as a plugin.
type Iglu interface {
	OnLoad() error
	PluginHTTP(request *http.Request) *http.Response
	GetManifest() PluginManifest
	RegisterDevice(reg DeviceRegistration) error
	OnDeviceToggle(id int, status bool) error
	GetWebExtensions() []WebExtension
}

// IgluRPC is what the server is using to communicate to the plugin over RPC
type IgluRPC struct{ client *rpc.Client }

type OnLoadReply struct {
	Err error
}

func (i *IgluRPC) OnLoad() error {
	rep := &OnLoadReply{}
	err := i.client.Call("Plugin.OnLoad", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Err
}

type Context struct {
	Username string
}

type PluginHTTPArgs struct {
	Request *http.Request
}

type PluginHTTPReply struct {
	Response *http.Response
}

func (i *IgluRPC) PluginHTTP(req *http.Request) *http.Response {
	args := &PluginHTTPArgs{Request: req}
	rep := &PluginHTTPReply{}
	err := i.client.Call("Plugin.PluginHTTP", args, &rep)
	if err != nil {
		panic(err)
	}
	return rep.Response
}

type GetManifestReply struct {
	Manifest PluginManifest
}

func (i *IgluRPC) GetManifest() PluginManifest {
	rep := &GetManifestReply{}
	err := i.client.Call("Plugin.GetManifest", new(interface{}), &rep)
	if err != nil {
		panic(err)
	}
	return rep.Manifest
}

type RegisterDeviceArgs struct {
	Reg DeviceRegistration
}

type RegisterDeviceReply struct {
	Err error
}

func (i *IgluRPC) RegisterDevice(reg DeviceRegistration) error {
	args := &RegisterDeviceArgs{Reg: reg}
	rep := &RegisterDeviceReply{}
	err := i.client.Call("Plugin.RegisterDevice", args, &rep)
	if err != nil {
		panic(err)
	}
	return rep.Err
}

type OnDeviceToggleArgs struct {
	Id     int
	Status bool
}

type OnDeviceToggleReply struct {
	Err error
}

func (i *IgluRPC) OnDeviceToggle(id int, status bool) error {
	args := &OnDeviceToggleArgs{Id: id, Status: status}
	rep := &OnDeviceToggleReply{}
	err := i.client.Call("Plugin.OnDeviceToggle", args, &rep)
	if err != nil {
		panic(err)
	}
	return rep.Err
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

// IgluRPCServer is the RPC server which IgluRPC talks to.
type IgluRPCServer struct {
	Impl Iglu
}

func (s *IgluRPCServer) OnLoad(args interface{}, resp *OnLoadReply) error {
	resp.Err = s.Impl.OnLoad()
	return nil
}

func (s *IgluRPCServer) PluginHTTP(args PluginHTTPArgs, resp *PluginHTTPReply) error {
	resp.Response = s.Impl.PluginHTTP(args.Request)
	return nil
}

func (s *IgluRPCServer) GetManifest(args interface{}, resp *GetManifestReply) error {
	resp.Manifest = s.Impl.GetManifest()
	return nil
}

func (s *IgluRPCServer) RegisterDevice(args RegisterDeviceArgs, resp *RegisterDeviceReply) error {
	resp.Err = s.Impl.RegisterDevice(args.Reg)
	return nil
}

func (s *IgluRPCServer) OnDeviceToggle(args OnDeviceToggleArgs, resp *OnDeviceToggleReply) error {
	resp.Err = s.Impl.OnDeviceToggle(args.Id, args.Status)
	return nil
}

func (s *IgluRPCServer) GetWebExtensions(args interface{}, resp *GetWebExtensionsReply) error {
	resp.Extensions = s.Impl.GetWebExtensions()
	return nil
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
