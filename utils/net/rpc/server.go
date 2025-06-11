package rpc;

import(
	IO       "io"
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	RPC      "net/rpc"
//	Time     "time"
	Sync     "sync"
	Atomic   "sync/atomic"
	TLS      "crypto/tls"
	Strings  "strings"
	Errors   "errors"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
);



type UpLink struct {
	MuxState   Sync.Mutex
	WaitGroup  *Sync.WaitGroup
	// transport
	Bind       string
	UseTLS     bool
	Listen     Net.Listener
	Server     *RPC.Server
	// stats
	Stats      *UpStats
	NextIndex  Atomic.Uint64
	NumReqs    Atomic.Uint64
	// state
	secret     string
	UserFunc   ValidUserFunc
	MuxSession Sync.Mutex
	Sessions   map[uint64]*Session
}

type Session struct {
	Conn Net.Conn
	User string
}

type ValidUserFunc func(string) bool;

type UpStats struct {
	CountConns uint64
	CountReqs  uint64
}



//TODO: Session[] cleanup
func NewUpLink(bind string) *UpLink {
	return &UpLink{
		Bind:     bind,
//		UseTLS:   true,
		Server:   RPC.NewServer(),
		secret:   "abcdefghijklmnopqrstuvwxyz",
		Sessions: make(map[uint64]*Session),
	};
}

func (link *UpLink) Start() error {
	link.MuxState.Lock();
	defer link.MuxState.Unlock();
	if link.Bind == "" { return Errors.New("Bind address is required"); }
	if link.secret != "" {
		if len(link.secret) <  8 { return Errors.New("Invalid secret length; too short"); }
		if len(link.secret) > 30 { return Errors.New("Invalid secret length; too long" ); }
	}
	if Strings.HasPrefix(link.Bind, "unix://") {
		link.UseTLS = false;
	}
	listen, err := UtilsNet.NewServerSocket(link.Bind);
	if err != nil {
		return Fmt.Errorf("%s%s for NewServerSocket in NewUpLink", LogPrefix, err);
	}
	link.Listen = listen;
	if link.WaitGroup == nil {
		var wait_group Sync.WaitGroup;
		link.WaitGroup = &wait_group;
	}
	go link.Serve();
	Utils.SleepC();
	return nil;
}

func (link *UpLink) Close() {
	link.MuxState.Lock();
	defer link.MuxState.Unlock();
	if link.Listen != nil {
		link.Listen.Close();
		link.Listen = nil;
	}
}

func (link *UpLink) CloseAll() {
	link.MuxState.Lock();
	defer link.MuxState.Unlock();
//TODO
}



func (link *UpLink) Serve() {
	link.WaitGroup.Add(1);
	defer link.WaitGroup.Done();
	Log.Printf("%sStarting RPC Server.. %s", LogPrefix, link.Bind);
	if link.UseTLS { Log.Printf("%sTLS Enabled",  LogPrefix);
	} else {         Log.Printf("%sTLS Disabled", LogPrefix); }
	config_tls := &TLS.Config{
		MinVersion: TLS.VersionTLS13,
//		Certificates: []TLS.Certificate{ cert },
//		ServerName: "pxn",
//		InsecureSkipVerify: true,
	};
	LOOP_SERVE:
	for {
		conn, err := link.Listen.Accept();
		if err != nil {
			Log.Printf("%s%s in UpLink->Serve()", LogPrefix, err);
			continue LOOP_SERVE;
		}
		if link.UseTLS {
			conn = TLS.Server(conn, config_tls);
		}
		go link.Handle(conn);
	}
}

func (link *UpLink) Handle(conn Net.Conn) {
	defer conn.Close();
	remote := conn.RemoteAddr().String();
	index := link.NextIndex.Add(1);
	Log.Printf("%s%d Connection from: %s", LogPrefix, index, remote);
	buffer := make([]byte, 64);
	if _, err := IO.ReadFull(conn, buffer); err != nil {
		Log.Printf("%s%s in UpLink->Handle()", LogPrefix, err);
		return;
	}
	gotsec := Strings.Trim(string(buffer[ 0:31]), "\000");
	gotusr := Strings.Trim(string(buffer[32:63]), "\000");
	if gotsec != link.secret {
		Log.Printf("%sInvalid secret from: %s", LogPrefix, remote);
		Utils.SleepR();
		conn.Close();
		return;
	}
	if link.UserFunc != nil {
		if !link.UserFunc(gotusr) {
			Log.Printf("%sInvalid user from: %s", LogPrefix, remote);
			Utils.SleepR();
			conn.Close();
			return;
		}
	}
	link.MuxSession.Lock();
	defer link.MuxSession.Unlock();
	session := Session{
		Conn: conn,
		User: gotusr,
	};
	link.Sessions[index] = &session;
	link.Server.ServeConn(conn);
}



//func (link *UpLink) GetStats() *UpStats {
//	type UpStats struct {
//		CountConns uint64
//		CountReqs  uint64
//	}
//}
