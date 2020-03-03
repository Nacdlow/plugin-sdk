package sdk

// IgluRPCServer is the RPC server which IgluRPC talks to.
type IgluRPCServer struct {
	Impl Iglu
}

func (s *IgluRPCServer) OnLoad(args interface{}, resp *OnLoadReply) error {
	resp.Err = s.Impl.OnLoad()
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

func (s *IgluRPCServer) GetPluginConfiguration(args interface{}, resp *GetPluginConfigurationReply) error {
	resp.Configuration = s.Impl.GetPluginConfiguration()
	return nil
}

func (s *IgluRPCServer) OnConfigurationUpdate(args OnConfigurationUpdateArgs, resp *interface{}) error {
	s.Impl.OnConfigurationUpdate(args.Configuration)
	return nil
}

func (s *IgluRPCServer) GetAvailableDevices(args interface{}, resp *GetAvailableDevicesReply) error {
	resp.Devices = s.Impl.GetAvailableDevices()
	return nil
}
