package rpc;

import(
//	Fmt      "fmt"
	Time     "time"
	Sync     "sync"
//	Errors   "errors"
	Utils    "github.com/PxnPub/pxnGoUtils/utils"
	UtilsNum "github.com/PxnPub/pxnGoUtils/utils/numbers"
	UtilsNet "github.com/PxnPub/pxnGoUtils/utils/net"
);



const DefaultParallel   = 128;
const DefaultReplyQueue =  "4k";
const DefaultBufferSize = "64k";
const DefaultFlushDelay = 0;



type FuncOnConnect func(remote string, request interface{}) (reply interface{});



type ServerRPC struct {
	WaitGroup  *Sync.WaitGroup
	StopGroup  *Sync.WaitGroup
	StopChan   chan struct{}
	MuxState   Sync.Mutex
	Bind       string
	Handler    *FuncHandler
//	Dispatch   *Dispatcher
	Parallel   uint16
	ReplyQueue uint16
	BufSizeOut uint16
	BufSizeIn  uint16
	FlushDelay Time.Duration
}



func NewServerRPC(bind string) (*ServerRPC, error) {
	listen, err := UtilsNet.NewServerSocket(bind);
	if err != nil { return nil, err; }





	def_replyqueue, err := UtilsNum.ParseByteSize(DefaultReplyQueue);
	if err != nil { return nil, err; }
	def_buffersize, err := UtilsNum.ParseByteSize(DefaultBufferSize);
	if err != nil { return nil, err; }


	return &ServerRPC{
		Bind:       bind,
		Parallel:   DefaultParallel,
		ReplyQueue: def_replyqueue,
		BufSizeOut: def_buffersize,
		BufSizeIn:  def_buffersize,
		FlushDelay: DefaultFlushDelay,
	};
}



func (server *ServerRPC) Start() error {
	server.MuxState.Lock();
	defer server.MuxState.Unlock();
//TODO: parameter checks
	if server.StopChan == nil {
		if err := server.Listen.Init(server.Bind); err != nil {
			return err;
		}
		workchans := make(chan struct{}, server.Parallel);
		go server.Serve(workchans);
		return nil;
	}
}



func (server *ServerRPC) Stop() {
	server.MuxState.Lock();
	defer server.MuxState.Unlock();
	if server.StopChan != nil {
		close(server.StopChan);
		server.StopGroup.Wait();
		server.StopChan = nil;
	}
}



func (server *ServerRPC) Serve(workchans chan struct{}) {
	if server.WaitGroup != nil {
		server.WaitGroup.Add(1);
		defer  server.WaitGroup.Done();
	}



}




















































}
