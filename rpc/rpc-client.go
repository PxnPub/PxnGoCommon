package rpc;

import(
	Log     "log"
	Fmt     "fmt"
	Time    "time"
	Sync    "sync"
	Context "context"
	Errors  "errors"
	GRPC    "google.golang.org/grpc"
	GInsec  "google.golang.org/grpc/credentials/insecure"
	GConty  "google.golang.org/grpc/connectivity"
//	GZIP    "google.golang.org/grpc/encoding/gzip"
	Service "github.com/PxnPub/PxnGoCommon/service"
	PxnNet  "github.com/PxnPub/PxnGoCommon/utils/net"
	PxnSan  "github.com/PxnPub/PxnGoCommon/utils/san"
	Utils   "github.com/PxnPub/PxnGoCommon/utils"
);



type ClientRPC struct {
	mut_state    Sync.Mutex
	service     *Service.Service
	// transport
	remote      string
	use_tls     bool
	grpc_client *GRPC.ClientConn
}



func NewClientRPC(service *Service.Service, remote string) *ClientRPC {
	return &ClientRPC{
		service: service,
		remote:  remote,
	};
}



func (rpc *ClientRPC) Start() error {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	if rpc.grpc_client != nil { return Errors.New("RPC client already started"); }
	if rpc.remote      == ""  { return Errors.New("RPC address is required"   ); }
	protocol, address, port := PxnNet.SplitProtocolAddressPort(rpc.remote);
	if protocol == "" { return Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		rpc.use_tls = false;
//TODO
panic("UNFINISHED UNIX RPC CLIENT");
		break;
	case "tcp", "tcp4", "tcp6":
		if rpc.use_tls { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !PxnSan.IsSafeDomain(address) {
			return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0 { return Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		backoff_maxdelay, err := Time.ParseDuration(DefaultBackoffMaxDelay);
		if err != nil { Log.Panicf("Invalid backoff max delay", err); }
		grpc_client, err := GRPC.NewClient(
			addrport,
			GRPC.WithTransportCredentials(GInsec.NewCredentials()),
			GRPC.WithBackoffMaxDelay(backoff_maxdelay),
		);
		if err != nil { return Fmt.Errorf("%s, failed to connect", err); }
		rpc.grpc_client = grpc_client;
		break;
	default: return Fmt.Errorf("Unknown protocol: %s", protocol);
	}
	go rpc.Serve();
	Utils.SleepC();
	return nil;
}

func (rpc *ClientRPC) Serve() {
	rpc.service.WaitGroup.Add(1);
	defer func() {
		rpc.Close();
		rpc.service.WaitGroup.Done();
	}();
//TODO: remove this
//	rpc.Service.AddCloseE(rpc);
	Log.Printf("%sConnecting RPC.. %s", LogPrefix, rpc.remote);
	rpc.grpc_client.WaitForStateChange(Context.Background(), GConty.Connecting);
	state := rpc.grpc_client.GetState();
	switch state {
	case GConty.Idle, GConty.Ready: break;
	default: Log.Panic("Connect state failure %s", state);
	}
//TODO
//TODO
//TODO
//TODO
//TODO: replace this with a health listener
//https://github.com/grpc/grpc-go/tree/v1.73.0/health
	last_state := state;
	LOOP_STATE:
	for {
		rpc.grpc_client.WaitForStateChange(Context.Background(), last_state);
		state := rpc.grpc_client.GetState();
		switch state {
		case GConty.Idle:       Log.Printf("%sIdle.. %s",       LogPrefix, rpc.remote);
		case GConty.Connecting: Log.Printf("%sConnecting.. %s", LogPrefix, rpc.remote);
		case GConty.Ready:      Log.Printf("%sReady. %s",       LogPrefix, rpc.remote);
		case GConty.TransientFailure:
			Log.Printf("%sReconnecting.. %s", LogPrefix, rpc.remote);
		case GConty.Shutdown: break LOOP_STATE;
		}
		last_state = state;
	}
//TODO
//TODO
//TODO
//TODO
//TODO
}



func (rpc *ClientRPC) Close() {
	rpc.service.WaitGroup.Add(1);
	rpc.mut_state.Lock();
	defer func() {
		rpc.mut_state.Unlock();
		rpc.service.WaitGroup.Done();
	}();
	if rpc.grpc_client != nil {
		if err := rpc.grpc_client.Close(); err != nil {
			Log.Printf("%v, in ClientRPC->Close()", err); }
		rpc.grpc_client = nil;
	}
}

func (rpc *ClientRPC) IsStopping() bool {
	return (rpc.grpc_client.GetState() == GConty.Shutdown);
}



func (rpc *ClientRPC) GetClientGRPC() *GRPC.ClientConn {
	return rpc.grpc_client;
}

func (rpc *ClientRPC) SetClientGRPC(grpc_client *GRPC.ClientConn) *GRPC.ClientConn {
	rpc.mut_state.Lock();
	defer rpc.mut_state.Unlock();
	previous := rpc.grpc_client;
	rpc.grpc_client = grpc_client;
	return previous;
}
