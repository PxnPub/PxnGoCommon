package net;

import(
	Net      "net"
	Errors   "errors"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
);



func NewClientSocket(remote string) (Net.Conn, error) {
return nil, nil;
//TODO
//	if remote == "" { return nil, Errors.New("remote address required"); }
//	protocol, address, port := SplitProtocolAddressPort(remote);
//	if protocol == "" { return nil, Errors.New("protocol is required"); }
//	switch protocol {
//	case "unix":
//		if len(address) < 5 { return nil, Fmt.Errorf("Invalid unix socket: %s", address ); }
//		resolved, err := Net.ResolveUnixAddr(protocol, address);
//		if err != nil { return nil, err; }
//		listen, err := Net.ListenUnix(protocol, resolved);
//		if err != nil { return nil, err; }
//		return listen, nil;
//	case "tcp4": fallthrough;
//	case "tcp6": fallthrough;
//	case "tcp":
//		if !UtilsSan.IsSafeDomain(address) { return nil, Fmt.Errorf("Invalid address: %s", address); }
//		if port == 0                       { return nil, Fmt.Errorf("Invalid port: %d"); }
//		addrport := Fmt.Sprintf("%s:%d", address, port);
//		resolved, err := Net.ResolveTCPAddr(protocol, addrport);
//		if err != nil { return nil, err; }
//		listen, err := Net.ListenTCP(protocol, resolved);
//		if err != nil { return nil, err; }
//		return listen, nil;
//TODO
//	case "tls": fallthrough;
//	case "ssl":
//		return nil, nil;
//	default: break;
//	}
//	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}



func NewClientUDP(remote string) (*Net.UDPConn, error) {
return nil, nil;
//TODO
//	if remote == "" { return nil, Errors.New("remote address required"); }
//	protocol, address, port := SplitProtocolAddressPort(remote);
//	if protocol == "" { return nil, Errors.New("protocol is required"); }
//	switch protocol {
//	case "udp4": fallthrough;
//	case "udp6": fallthrough;
//	case "udp":
//		if !UtilsSan.IsSafeDomain(address) { return nil, Fmt.Errorf("Invalid address: %s", address); }
//		if port == 0                       { return nil, Fmt.Errorf("Invalid port: %d"); }
//		addrport := Fmt.Sprintf("%s:%d", address, port);
//		resolved, err := Net.ResolveUDPAddr(protocol, addrport);
//		if err != nil { return nil, err; }
//		listen, err := Net.ListenUDP(protocol, resolved);
//		if err != nil { return nil, err; }
//		return listen, nil;
//	default: break;
//	}
//	return nil, Fmt.Errorf("Unknown protocol: %s", protocol);
}
