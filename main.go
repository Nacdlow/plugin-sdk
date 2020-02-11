package api

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	macaron "gopkg.in/macaron.v1"
)

// Iglu is the interface that we're exposing as a plugin.
type Iglu interface {
	OnLoad()
	Middleware() macaron.Handler
	// TODO these are interfaces yet to be implemented
	RegisterDevice()
	OnDeviceToggle()
}

type IgluRPC struct{ client *rpc.Client }

func (i *IgluRPC) OnLoad() {
	err := i.client.Call("Plugin.OnLoad", new(interface{}), nil)
	if err != nil {
		panic(err)
	}
}

func (i *IgluRPC) Middleware() (handler macaron.Handler) {
	err := i.client.Call("Plugin.Middleware", new(interface{}), &handler)
	if err != nil {
		panic(err)
	}
	return
}

type IgluRPCServer struct {
	Impl Iglu
}

func (s *IgluRPCServer) OnLoad(args interface{}) error {
	s.Impl.OnLoad()
	return nil
}

func (s *IgluRPCServer) Middleware(args interface{}) macaron.Handler {
	return s.Impl.Middleware()
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
