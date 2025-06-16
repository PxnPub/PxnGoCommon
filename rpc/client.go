package rpc;

import(
	Log      "log"
	Fmt      "fmt"
	Sync     "sync"
	Errors   "errors"
	GRPC     "google.golang.org/grpc"
	GRPC_Ins "google.golang.org/grpc/credentials/insecure"
	_        "google.golang.org/grpc/encoding/gzip"
	Service  "github.com/PxnPub/PxnGoCommon/service"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
);



type Client struct {
	MuxState   Sync.Mutex
	Service    *Service.Service
	// transport
	Remote     string
	UseTLS     bool
	RPC        GRPC.ClientConnInterface
}



func NewClient(service *Service.Service, remote string) *Client {
	return &Client{
		Service: service,
		Remote:  remote,
	};
}



func (client *Client) Start() error {
	client.MuxState.Lock();
	defer client.MuxState.Unlock();
	if client.Remote == "" { return Errors.New("Broker address is required"); }
	Log.Printf("%sConnecting RPC.. %s", LogPrefix, client.Remote);
	protocol, address, port := UtilsNet.SplitProtocolAddressPort(client.Remote);
	if protocol == "" { return Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		client.UseTLS = false;
//TODO
panic("UNFINISHED UNIX RPC CLIENT");
		break;
	case "tcp", "tcp4", "tcp6":
		if client.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		if !UtilsSan.IsSafeDomain(address) { return Fmt.Errorf("Invalid address: %s", address); }
		if port == 0                       { return Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		rpc, err := GRPC.NewClient(
			addrport,
			GRPC.WithTransportCredentials(GRPC_Ins.NewCredentials()),
		);
		if err != nil { return Fmt.Errorf("%s failed to connect", err); }
		client.RPC = rpc;
		Utils.SleepC();
		return nil;
	default: break;
	}
	return Fmt.Errorf("Unknown protocol: %s", protocol);
}

func (client *Client) Close() {
//TODO
}
