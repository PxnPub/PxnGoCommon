package rpc;

import(
);



type Server struct {
	MuxState Sync.Mutex
	Service  *Service.Service
	// transport
	Bind   string
	UseTLS bool
	Listen Net.Listener
}



func NewServer(service *Service.Service, bind string) *Server {
	return &Server{
		Service: service,
		Bind:    bind,
	};
}



func (server *Server) Start() error {
	server.MuxState.Lock();
	defer server.MuxState.Unlock();
	if server.Bind == "" { return Errors.New("Bind address is required"); }
	Log.Printf("%sStarting RPC Server.. %s", LogPrefix, server.Bind);
	protocol, address, port := UtilsNet.SplitProtocolAddressPort(server.Remote);
	if protocol == "" { return Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		server.UseTLS = false;
//TODO
panic("UNFINISHED UNIX RPC SERVER");
		break;
	case "tcp", "tcp4", "tcp6":
		if server.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !UtilsSan.IsSafeDomain(address) { return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0                       { return Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		listen, err := UtilsNet.NewServerSocket(server.Bind);
		if err != nil { return Fmt.Errorf("%s failed to listen", err); }
		server.Listen = listen;
		server.Service.AddStopHook(func() {
			server.Listen.Close();
		});
		server.RPC = GRPC.NewServer();
print("AAAAA\n");
server.RPC.Serve(listen);
print("BBBBB\n");
		Utils.SleepC();
		return nil;
	default: break;
	}
	return Fmt.Errorf("Unknown protocol: %s", protocol);
}

func (server *Server) Close() {
//TODO
}
