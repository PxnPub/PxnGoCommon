package net;

import(
	OS       "os"
	Fmt      "fmt"
	Strings  "strings"
	StrConv  "strconv"
	UtilsSan "github.com/PxnPub/PxnGoCommon/utils/san"
);



func SplitProtocolAddressPort(bind string) (string, string, uint16) {
	var protocol string;
	var address  string = bind;
	var port     uint16;
	if Strings.Contains(address, "://") {
		parts := Strings.SplitN(address, "://", 2);
		protocol = parts[0];
		address  = parts[1];
	}
	if Strings.Contains(address, ":") {
		parts := Strings.SplitN(address, ":", 2);
		address = parts[0];
		p, err := StrConv.Atoi(parts[1]);
		if err != nil { panic(err); }
		port = uint16(p);

	}
	return protocol, address, port;
}



func RemoveOldUnixSocket(file string) error {
	if !UtilsSan.IsSafeFilePath(file) { return Fmt.Errorf("Invalid address: %s", file); }
	// file exists
	if _, err := OS.Stat(file); err == nil {
		// file type
		info, err := OS.Lstat(file);
		if err != nil { return Fmt.Errorf("Failed to stat file type: %v", err); }
		// is a socket
		if info.Mode()&OS.ModeSocket != 0 {
			// remove old socket file
			if err := OS.Remove(file); err != nil {
				return Fmt.Errorf("Failed to remove old socket file: %v", err);
			}
			return Fmt.Errorf("Removed old socket file: %s", file);
		}
		return Fmt.Errorf("File exists but is not a socket: %s", file);
	}
	return nil;
}
