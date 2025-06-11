package rpc;

import(
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	RPC      "net/rpc"
	Sync     "sync"
	Atomic   "sync/atomic"
	TLS      "crypto/tls"
	Strings  "strings"
	Errors   "errors"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
);



type BackLink struct {
	MuxState   Sync.Mutex
	WaitGroup  *Sync.WaitGroup
	// transport
	RemoteAddr string
	UseTLS     bool
	Client     *RPC.Client
	// stats
	Stats      *BackStats
	NextIndex  Atomic.Uint64
	NumReqs    Atomic.Uint64
	// state
	secret     string
	user       string
}

type BackStats struct {
	CountConns uint64
	CountReqs  uint64
}



func NewBackLink(addr string) *BackLink {
	return &BackLink{
		RemoteAddr: addr,
//		UseTLS:     true,
		secret:     "abcdefghijklmnopqrstuvwxyz",
		user:       "lop",
	};
}

func (link *BackLink) Start() error {
	link.MuxState.Lock();
	defer link.MuxState.Unlock();
	if link.RemoteAddr == "" { return Errors.New("Remote address is required"); }
	if link.secret != "" {
		if len(link.secret) <  8 { return Errors.New("Invalid secret length; too short"); }
		if len(link.secret) > 30 { return Errors.New("Invalid secret length; too long" ); }
	}
	if Strings.HasPrefix(link.RemoteAddr, "unix://") {
		link.UseTLS = false;
	}
	if link.WaitGroup == nil {
		var wait_group Sync.WaitGroup;
		link.WaitGroup = &wait_group;
	}
	if err := link.ConnectLoop(); err != nil { return err; }
	Utils.SleepC();
	return nil;
}

func (link *BackLink) Close() {
	link.MuxState.Lock();
	defer link.MuxState.Unlock();
print("CLOSE\n");
	if link.Client != nil {
		link.Client.Close();
		link.Client = nil;
	}
}



func (link *BackLink) ConnectLoop() error {
//TODO: reconnect loop
//	for {
		Log.Printf("Connecting to RPC.. %s", link.RemoteAddr);
		if link.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
		} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
		conn, err := UtilsNet.NewClientSocket(link.RemoteAddr);
		if err != nil { return Fmt.Errorf(
			"%s%s for NewClientSocket in NewBackLink",
			LogPrefix, err); }
		go link.Handle(conn);
		return nil;
//	}
}

func (link *BackLink) Handle(conn Net.Conn) {
	link.WaitGroup.Add(1);
	defer link.WaitGroup.Done();
	config_tls := &TLS.Config{
		MinVersion: TLS.VersionTLS13,
//		Certificates: []TLS.Certificate{ cert },
//		ServerName: "pxn",
//		InsecureSkipVerify: true,
	};
	if link.UseTLS {
		conn = TLS.Client(conn, config_tls);
	}
	buffer := make([]byte, 64);
	for i:=0; i<64; i++ { buffer[i] = 0x0; }
	copy(buffer[ 1:], []byte(link.secret));
	copy(buffer[32:], []byte(link.user));
	if _, err := conn.Write(buffer); err != nil {
		Log.Panicf("%s%s when writing secret to: %s", LogPrefix, err, link.RemoteAddr);
	}
	Utils.SleepC();
	client := RPC.NewClient(conn);
	link.Client = client;

}



func (link *BackLink) Call(name string) ([]byte, error) {
Fmt.Printf("CALL: %s\n", name);
return []byte("{}"), nil;
}



//func (link *BackLink) GetStats() *BackStats {
//	type BackStats struct {
//		CountConns uint64
//		CountReqs  uint64
//	}
//}
