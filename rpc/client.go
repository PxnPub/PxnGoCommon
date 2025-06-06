package rpc;

import(
//	Fmt      "fmt"
	Errors   "errors"
	UtilsNet "github.com/PxnPub/pxnGoUtils/utils/net"
);



type ClientRPC struct {
	Dispatch *Dispatcher
}



func NewClientRPC(remote string) (*ClientRPC, error) {
	listen, err := UtilsNet.NewClientSocket(remote);
	if err != nil { return nil, err; }












	return &ClientRPC{};
}
