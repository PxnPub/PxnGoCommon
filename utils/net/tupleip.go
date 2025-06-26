package net;

import(
	Fmt     "fmt"
	Net     "net"
	Strings "strings"
	StrConv "strconv"
	Binary  "encoding/binary"
);



type TupIP struct {
	H uint64
	L uint64
}



func ParseAddrStr(addr string) *TupIP {
	var ip Net.IP = Net.ParseIP(addr);
	if ip == nil { return nil; }
	// ipv4
	if ip.To4() != nil {
		ip4 := ip.To4();
		return &TupIP{
			H: 0,
			L: uint64(Binary.BigEndian.Uint32(ip4)),
		};
	// ipv6
	} else {
		ip6 := ip.To16();
		return &TupIP{
			H: Binary.BigEndian.Uint64(ip6[0: 8]),
			L: Binary.BigEndian.Uint64(ip6[8:16]),
		};
	}
}

func ParseTupStr(tup string) *TupIP {
	if tup == "" { return nil; }
	parts := Strings.SplitN(tup, ";", 2);
	if len(parts) != 2 { return nil; }
	ip_h, err := StrConv.ParseUint(parts[0], 10, 64);
	if err != nil { return nil; }
	ip_l, err := StrConv.ParseUint(parts[1], 10, 64);
	if err != nil { return nil; }
	return &TupIP{
		H: ip_h,
		L: ip_l,
	};
}

func (tup *TupIP) String() string {
	return Fmt.Sprintf("%d;%d", tup.H, tup.L);
}

//func (tup *TupIP) ToStringReal() string {
//TODO: should this function even exist?
//}
