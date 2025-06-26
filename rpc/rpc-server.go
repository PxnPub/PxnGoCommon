package rpc;

import(
	Log     "log"
	Fmt     "fmt"
	Net     "net"
	Sync    "sync"
	Errors  "errors"
	GRPC    "google.golang.org/grpc"
	_       "google.golang.org/grpc/encoding/gzip"
	Service "github.com/PxnPub/PxnGoCommon/service"
	PxnNet  "github.com/PxnPub/PxnGoCommon/utils/net"
	PxnSan  "github.com/PxnPub/PxnGoCommon/utils/san"
	Utils   "github.com/PxnPub/PxnGoCommon/utils"
);



type ServerRPC struct {
	mut_state   Sync.Mutex
	service     *Service.Service
	// transport
	bind        string
	use_tls     bool
	listen      Net.Listener
	grpc_server *GRPC.Server
}



func NewServerRPC(service *Service.Service, bind string) *ServerRPC {
	return &ServerRPC{
		service: service,
		bind:    bind,
	};
}



func (rpc *ServerRPC) Start() error {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.bind == "" { rpc.bind = DefaultBindRPC; }
	if rpc.bind == "" { return Errors.New("Bind address is required"); }
	protocol, address, port := PxnNet.SplitProtocolAddressPort(rpc.bind);
	if protocol == "" { return Errors.New("protocol is required"); }
	Log.Printf("Starting RPC Server.. %s", rpc.bind);
	if rpc.grpc_server == nil { rpc.grpc_server = GRPC.NewServer(); }
	switch protocol {
	case "unix":
		rpc.use_tls = false;
//TODO
panic("UNFINISHED UNIX RPC SERVER");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.use_tls { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !PxnSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d"); }
		listen, err := PxnNet.NewServerSocket(rpc.bind);
		if err != nil { return Fmt.Errorf("%s, failed to listen", err); }
		rpc.listen = listen;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s", protocol);
	}
	go rpc.Serve();
	Utils.SleepC();
	return nil;
}

func (rpc *ServerRPC) Serve() {
	rpc.service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.service.WaitGroup.Done();
	}();
	rpc.service.AddClose(rpc);
	if err := rpc.grpc_server.Serve(rpc.listen); err != nil {
		Log.Printf("%s in ServerRPC->Serve()", err); }
}



func (rpc *ServerRPC) Close() {
	rpc.service.WaitGroup.Add(1);
	defer rpc.service.WaitGroup.Done();
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.listen != nil {
		if err := rpc.listen.Close(); err != nil {
			Log.Printf("%v, in ServerRPC->Close()", err); }
		rpc.grpc_server.GracefulStop();
		rpc.listen = nil;
	}
}



func (rpc *ServerRPC) GetServerGRPC() *GRPC.Server {
	return rpc.grpc_server;
}

func (rpc *ServerRPC) SetServerGRPC(grpc_server *GRPC.Server) *GRPC.Server {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	previous := rpc.grpc_server;
	rpc.grpc_server = grpc_server;
	return previous;
}
