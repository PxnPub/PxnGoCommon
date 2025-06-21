package rpc;

import(
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	Sync     "sync"
	Errors   "errors"
	GRPC     "google.golang.org/grpc"
	_        "google.golang.org/grpc/encoding/gzip"
	Service  "github.com/PxnPub/PxnGoCommon/service"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
);



type RPCServer struct {
	MutState Sync.Mutex
	Service  *Service.Service
	// transport
	Bind     string
	UseTLS   bool
	Listen   Net.Listener
	Server   *GRPC.Server
}



func NewRPCServer(service *Service.Service, bind string) *RPCServer {
	return &RPCServer{
		Service: service,
		Bind:    bind,
	};
}



func (rpc *RPCServer) Start() error {
	rpc.MutState.Lock();
	defer rpc.MutState.Unlock();
	if rpc.Bind == "" { rpc.Bind = DefaultBindRPC; }
	if rpc.Bind == "" { return Errors.New("Bind address is required"); }
	protocol, address, port := UtilsNet.SplitProtocolAddressPort(rpc.Bind);
	if protocol == "" { return Errors.New("protocol is required"); }
	Log.Printf("Starting RPC Server.. %s", rpc.Bind);
	if rpc.Server == nil { rpc.Server = GRPC.NewServer(); }
	switch protocol {
	case "unix":
		rpc.UseTLS = false;
//TODO
panic("UNFINISHED UNIX RPC SERVER");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {        Log.Printf("%sTLS Disabled", LogPrefix); }
		if !UtilsSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d"); }
		listen, err := UtilsNet.NewServerSocket(rpc.Bind);
		if err != nil { return Fmt.Errorf("%s, failed to listen", err); }
		rpc.Listen = listen;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s", protocol);
	}
	go rpc.Serve();
	Utils.SleepC();
	return nil;
}

func (rpc *RPCServer) Serve() {
	rpc.Service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.Service.WaitGroup.Done();
	}();
	rpc.Service.AddCloseE(rpc);
	if err := rpc.Server.Serve(rpc.Listen); err != nil {
		Log.Printf("%s in RPCServer->Serve()", err); }
}



func (rpc *RPCServer) Close() error {
	rpc.Service.WaitGroup.Add(1);
	defer rpc.Service.WaitGroup.Done();
	rpc.MutState.Lock();
	defer rpc.MutState.Unlock();
	var e error = nil;
	if rpc.Listen != nil {
		e = rpc.Listen.Close();
		rpc.Listen = nil;
	}
	rpc.Server.GracefulStop();
	return e;
}
