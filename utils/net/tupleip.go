package net;

import(
	Net    "net"
	Binary "encoding/binary"
);



type TupleIP struct {
	H uint64
	L uint64
}



func StringToIntIP(address string) (*TupleIP, error) {
	var ip Net.IP = Net.ParseIP(address);
	if ip == nil { return nil, Fmt.Errorf("Invalid address: %s", address); }
	// ipv4
	if ip.To4() != nil {
		ip4 := ip.To4();
		return &TupleIP{
			H: 0,
			L: uint64(Binary.BigEndian.Uint32(ip4)),
		}, nil;
	// ipv6
	} else {
		ip6 := ip.To16();
		return &TupleIP{
			H: Binary.BigEndian.Uint64(ip6[0:8]),
			L: Binary.BigEndian.Uint64(ip6[8:16]),
		}, nil;
	}
}
