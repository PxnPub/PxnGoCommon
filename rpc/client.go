package rpc;

import(
	Log      "log"
	Fmt      "fmt"
	Sync     "sync"
	Context  "context"
	Errors   "errors"
	GRPC     "google.golang.org/grpc"
	GRPC_Ins "google.golang.org/grpc/credentials/insecure"
	GConty   "google.golang.org/grpc/connectivity"
//	GZIP     "google.golang.org/grpc/encoding/gzip"
	Service  "github.com/PxnPub/PxnGoCommon/service"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
);



type RPCClient struct {
	MutState Sync.Mutex
	Service  *Service.Service
	// transport
	Remote   string
	UseTLS   bool
	Client   *GRPC.ClientConn
}



func NewRPCClient(service *Service.Service, remote string) *RPCClient {
	return &RPCClient{
		Service: service,
		Remote:  remote,
	};
}



func (rpc *RPCClient) Start() error {
	rpc.MutState.Lock();
	defer rpc.MutState.Unlock();
	if rpc.Client != nil { return Errors.New("RPC client already started"); }
	if rpc.Remote == ""  { return Errors.New("RPC address is required"   ); }
	protocol, address, port := UtilsNet.SplitProtocolAddressPort(rpc.Remote);
	if protocol == "" { return Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		rpc.UseTLS = false;
//TODO
panic("UNFINISHED UNIX RPC CLIENT");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {        Log.Printf("%sTLS Disabled", LogPrefix); }
		if !UtilsSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		client, err := GRPC.NewClient(
			addrport,
			GRPC.WithTransportCredentials(GRPC_Ins.NewCredentials()),
		);
		if err != nil { return Fmt.Errorf("%s, failed to connect", err); }
		rpc.Client = client;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s", protocol);
	}
	go rpc.Serve();
	Utils.SleepC();
	return nil;
}

func (rpc *RPCClient) Serve() {
	rpc.Service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.Service.WaitGroup.Done();
	}();
	rpc.Service.AddCloseE(rpc);
	Log.Printf("%sConnecting RPC.. %s", LogPrefix, rpc.Remote);
	rpc.Client.WaitForStateChange(Context.Background(), GConty.Connecting);
	state := rpc.Client.GetState();
	switch state {
	case GConty.Idle, GConty.Ready: break;
	default: Log.Panic("Connect state failure %s", state);
	}
//TODO: replace this with a health listener
//https://github.com/grpc/grpc-go/tree/v1.73.0/health
	last_state := state;
	LOOP_STATE:
	for {
		rpc.Client.WaitForStateChange(Context.Background(), last_state);
		state := rpc.Client.GetState();
		switch state {
		case GConty.Idle:             Log.Printf("%sIdle.. %s",         LogPrefix, rpc.Remote);
		case GConty.Connecting:       Log.Printf("%sConnecting.. %s",   LogPrefix, rpc.Remote);
		case GConty.Ready:            Log.Printf("%sReady. %s",         LogPrefix, rpc.Remote);
		case GConty.TransientFailure: Log.Printf("%sReconnecting.. %s", LogPrefix, rpc.Remote);
		case GConty.Shutdown:         break LOOP_STATE;
		}
		last_state = state;
	}
}



func (rpc *RPCClient) Close() error {
	rpc.Service.WaitGroup.Add(1);
	defer rpc.Service.WaitGroup.Done();
	rpc.MutState.Lock();
	defer rpc.MutState.Unlock();
	var e error = nil;
	if rpc.Client != nil {
		if err := rpc.Client.Close(); err != nil { e = err; }}
	return e;
}
