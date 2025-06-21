package web;

import(
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	HTTP     "net/http"
	Time     "time"
	Sync     "sync"
//	Atomic   "sync/atomic"
	Context  "context"
	Errors   "errors"
	Gorilla  "github.com/gorilla/mux"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
	Service  "github.com/PxnPub/PxnGoCommon/service"
);



const LogPrefix      = "[Web] ";
const DefaultBindWeb = "tcp://127.0.0.1:8000";



type WebServer struct {
	MutState Sync.Mutex
	Service  *Service.Service
	// transport
	Bind   string
	UseTLS bool
	Listen Net.Listener
	Router *Gorilla.Router
	Server *HTTP.Server
}



func NewWebServer(service *Service.Service, bind string) *WebServer {
	web := WebServer{
		Service: service,
		Bind:    bind,
		Router:  Gorilla.NewRouter(),
		Server:  &HTTP.Server{},
	};
	web.Router.NotFoundHandler = HTTP.HandlerFunc(web.PageNotFound);
	web.Router.Use(web.MiddlewareStats);
	return &web;
}



func (web *WebServer) Start() error {
	web.MutState.Lock();
	defer web.MutState.Unlock();
	if web.Bind == "" { web.Bind = DefaultBindWeb; }
	if web.Bind == "" { return Errors.New("Bind address is required"); }
	listen, err := UtilsNet.NewServerSocket(web.Bind);
	if err != nil { return Fmt.Errorf(
		"%s for NewServerSocket() in WebServer->Start()", err); }
	web.Listen = listen;
	go web.Serve();
	Utils.SleepC();
	return nil;
}

func (web *WebServer) Serve() {
	web.Service.WaitGroup.Add(1);
	defer func() {
		web.Close();
		web.Service.WaitGroup.Done();
	}();
	web.Service.AddCloseE(web);
	Log.Printf("Starting Web Server.. %s", web.Bind);
	web.Server.Handler = web.Router;
	if err := web.Server.Serve(web.Listen); err != nil {
		Log.Printf("%s, in WebServer->Serve()", err); }
}



func (web *WebServer) Close() error {
	web.Service.WaitGroup.Add(1);
	defer web.Service.WaitGroup.Done();
	web.MutState.Lock();
	defer web.MutState.Unlock();
	var e error = nil;
	if web.Listen != nil {
		if err := web.Listen.Close(); err != nil { e = err; }
		web.Listen = nil;
	}
	if err := web.Server.Shutdown(Context.Background()); err != nil { e = err; }
	return e;
}



func (web *WebServer) MiddlewareStats(next HTTP.Handler) HTTP.Handler {
	return HTTP.HandlerFunc(func(w HTTP.ResponseWriter, r *HTTP.Request) {
		start := Time.Now();
		next.ServeHTTP(w, r);
		Log.Printf("%s%s %s in %v", LogPrefix, r.Method, r.URL.Path, Time.Since(start));
	});
}

func (web *WebServer) PageNotFound(w HTTP.ResponseWriter, r *HTTP.Request) {
	HTTP.Error(w, "404 Not Found", HTTP.StatusNotFound);
	Log.Printf("%s404 %s %s", LogPrefix, r.Method, r.URL.Path);
}



func AddStaticRoute(router *Gorilla.Router) {
	fs := HTTP.FileServer(HTTP.Dir("./static"));
	router.PathPrefix("/static/").Handler(HTTP.StripPrefix("/static/", fs));
}



//func (web *WebServer) GetStats() *Stats {
//	type StatsRPC struct {
//		CountConns uint64
//		CountReqs  uint64
//	}
//}



func NewRedirect(target string) HTTP.HandlerFunc  {
	return func(out HTTP.ResponseWriter, in *HTTP.Request) {
		HTTP.Redirect(out, in, target, HTTP.StatusFound);
	};
}
