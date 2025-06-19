package rpc;

import(
	Log     "log"
	Fmt     "fmt"
	Net     "net"
	Context "context"
	GRPC    "google.golang.org/grpc"
	GPeer   "google.golang.org/grpc/peer"
	GStatus "google.golang.org/grpc/status"
	GCodes  "google.golang.org/grpc/codes"
);



const KeyUsername = "username";



func NewAuthByIP(users map[string]string) GRPC.UnaryServerInterceptor {
	return func(ctx Context.Context, req any, info *GRPC.UnaryServerInfo,
			handler GRPC.UnaryHandler) (any, error) {
		if _, ok := ctx.Value(KeyUsername).(int); ok {
			return handler(ctx, req);
		}
		peer, ok := GPeer.FromContext(ctx);
		if !ok {
			Log.Printf("%sUnable to find peer info", LogPrefix);
			return nil, GStatus.Errorf(
				GCodes.PermissionDenied,
				"Unable to find peer info",
			);
		}
		remote := peer.Addr.String();
		if addr, _, err := Net.SplitHostPort(remote); err == nil {
			if addr != "" { remote = addr; }}
		if username, ok := users[remote]; ok {
			ctx = Context.WithValue(ctx, KeyUsername, username);
			return handler(ctx, req);
		}
		Log.Printf("%sIP not allowed: %s", LogPrefix, remote);
		return nil, GStatus.Errorf(
			GCodes.PermissionDenied,
			Fmt.Sprintf("IP not allowed: %s", remote),
		);
	};
}
