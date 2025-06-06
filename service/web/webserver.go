package web;

import(
	Log      "log"
	Time     "time"
	HTTP     "net/http"
	Context  "context"
	Atomic   "sync/atomic"
	Gorilla  "github.com/gorilla/mux"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
	Service  "github.com/PxnPub/PxnGoCommon/service"
);



const DefaultBindWeb = "tcp://127.0.0.1:8000";



type WebServer struct {
	TrapC   *Service.TrapC
	Server  *HTTP.Server
	Router  *Gorilla.Router
	Bind    string
	StatReq Atomic.Uint64
}



func NewWebServer(trapc *Service.TrapC, bind string) *WebServer {
	web := WebServer{
		TrapC:  trapc,
		Router: Gorilla.NewRouter(),
		Bind:   bind,
	};
	web.Router.Use(MiddlewareStats(&web));
	return &web;
}



func (web *WebServer) Start() error {
	listen, err := UtilsNet.NewListenerSocket(web.Bind);
	if err != nil { return err; }
	go func () {
		web.TrapC.WaitGroup.Add(1);
		defer web.TrapC.WaitGroup.Done();
		web.TrapC.AddStopHook(func() {
			if err := web.Server.Shutdown(Context.Background()); err != nil {
				Log.Printf("HTTP Shutdown Error: %v", err);
			}
		});
		Log.Printf("[%s] Listening..", web.Bind);
		web.Server = &HTTP.Server{ Handler: web.Router };
		if err := web.Server.Serve(listen); err != nil {
			CASE_ERR:
			switch err.Error() {
			case "http: Server closed":
				Log.Printf("[%s] Listener closed.", web.Bind);
				break CASE_ERR;
			default: Log.Printf("HTTP Listen Error: %v", err);
			}
		}
	}();
	web.Close();
	Utils.SleepC();
	return nil;
}



func (web *WebServer) Close() {
	if server := web.Server; server != nil {
		server.Close();
	}
}



func MiddlewareStats(web *WebServer) Gorilla.MiddlewareFunc {
	return func(next HTTP.Handler) HTTP.Handler {
		return HTTP.HandlerFunc(func(out HTTP.ResponseWriter, in *HTTP.Request) {
			cnt := web.StatReq.Add(1);
			Log.Printf("REQUEST[%d] %s\n", cnt, in.RequestURI);
			next.ServeHTTP(out, in);
		});
	}
}



func AddRouteStatic(router *Gorilla.Router) {
	fs := HTTP.FileServer(HTTP.Dir("./static"));
	router.PathPrefix("/static/").Handler(HTTP.StripPrefix("/static/", fs));
}
