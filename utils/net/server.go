package net;

import(
	Fmt      "fmt"
	Net      "net"
	Errors   "errors"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
);



func NewServerSocket(bind string) (Net.Listener, error) {
	if bind == "" { return nil, Errors.New("bind address required"); }
	protocol, address, port := SplitProtocolAddressPort(bind);
	if protocol == "" { return nil, Errors.New("protocol is required"); }
	switch protocol {
	case "unix":
		if len(address) < 5 { return nil, Fmt.Errorf("Invalid unix socket: %s", address ); }
		if err := RemoveOldUnixSocket(address); err != nil { return nil, err; }
		resolved, err := Net.ResolveUnixAddr(protocol, address);
		if err != nil { return nil, err; }
//TODO: is this right?
		listen, err := Net.ListenUnix(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
	case "tcp4": fallthrough;
	case "tcp6": fallthrough;
	case "tcp":
		if !UtilsSan.IsSafeDomain(address) { return nil, Fmt.Errorf("Invalid address: %s", address); }
		if port == 0                       { return nil, Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		resolved, err := Net.ResolveTCPAddr(protocol, addrport);
		if err != nil { return nil, err; }
		listen, err := Net.ListenTCP(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
//TODO
//	case "tls": fallthrough;
//	case "ssl":
//		return nil, nil;
	default: break;
	}
	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}



func NewServerUDP(bind string) (*Net.UDPConn, error) {
	if bind == "" { return nil, Errors.New("bind address required"); }
	protocol, address, port := SplitProtocolAddressPort(bind);
	if protocol == "" { return nil, Errors.New("protocol is required"); }
	switch protocol {
	case "udp4": fallthrough;
	case "udp6": fallthrough;
	case "udp":
		if !UtilsSan.IsSafeDomain(address) { return nil, Fmt.Errorf("Invalid address: %s", address); }
		if port == 0                       { return nil, Fmt.Errorf("Invalid port: %d"); }
		addrport := Fmt.Sprintf("%s:%d", address, port);
		resolved, err := Net.ResolveUDPAddr(protocol, addrport);
		if err != nil { return nil, err; }
		listen, err := Net.ListenUDP(protocol, resolved);
		if err != nil { return nil, err; }
		return listen, nil;
	default: break;
	}
	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}
