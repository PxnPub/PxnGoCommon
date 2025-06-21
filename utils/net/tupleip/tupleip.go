package tupleip;

import(
	Fmt     "fmt"
	Net     "net"
	Strings "strings"
	StrConv "strconv"
	Binary  "encoding/binary"
);



type IP struct {
	H uint64
	L uint64
}



func NewFromString(address string) (*IP, error) {
	var ip Net.IP = Net.ParseIP(address);
	if ip == nil { return nil, Fmt.Errorf("Invalid address: %s", address); }
	// ipv4
	if ip.To4() != nil {
		ip4 := ip.To4();
		return &IP{
			H: 0,
			L: uint64(Binary.BigEndian.Uint32(ip4)),
		}, nil;
	// ipv6
	} else {
		ip6 := ip.To16();
		return &IP{
			H: Binary.BigEndian.Uint64(ip6[0: 8]),
			L: Binary.BigEndian.Uint64(ip6[8:16]),
		}, nil;
	}
}

func Parse(ip string) *IP {
	if ip == "" { return nil; }
	parts := Strings.SplitN(ip, ";", 2);
	if len(parts) != 2 { return nil; }
	ip_h, err := StrConv.ParseUint(parts[0], 10, 64);
	if err != nil { return nil; }
	ip_l, err := StrConv.ParseUint(parts[1], 10, 64);
	if err != nil { return nil; }
	return &IP{
		H: ip_h,
		L: ip_l,
	};
}



func (ip *IP) ToStringRaw() string {
	return Fmt.Sprintf("%d;%d", ip.H, ip.L);
}

//func (ip *IP) ToStringReal() string {
//TODO: should this function even exist?
//}
