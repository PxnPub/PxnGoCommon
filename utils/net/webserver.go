package net;

import(
	Log     "log"
	Fmt     "fmt"
	Net     "net"
	HTTP    "net/http"
	Time    "time"
	Sync    "sync"
	Strings "strings"
	Context "context"
	Errors  "errors"
	Gorilla "github.com/gorilla/mux"
	Utils   "github.com/PxnPub/PxnGoCommon/utils"
	Service "github.com/PxnPub/PxnGoCommon/service"
);



const LogPrefixWeb   = "[Web] ";
const DefaultBindWeb = "tcp://127.0.0.1:8000";



type WebServer struct {
	mut_state Sync.Mutex
	service   *Service.Service
	// transport
	bind      string
	use_tls   bool
	listen    Net.Listener
	Router    HTTP.Handler
	server    *HTTP.Server
}



func NewWebServer(service *Service.Service, bind string) *WebServer {
	web := WebServer{
		service: service,
		bind:    bind,
		server:  &HTTP.Server{},
	};
	return &web;
}



func (web *WebServer) WithHandler(router HTTP.Handler) *WebServer {
	web.Router = router;
	return web;
}

func (web *WebServer) WithGorilla() *Gorilla.Router {
	router := Gorilla.NewRouter();
	router.NotFoundHandler = HTTP.HandlerFunc(web.PageNotFound);
	router.Use(web.MiddlewareStats);
	web.Router = router;
	return router;
}



func (web *WebServer) Start() error {
	web.mut_state.Lock();
	defer web.mut_state.Unlock();
	if web.bind == "" { web.bind = DefaultBindWeb; }
	if web.bind == "" { return Errors.New("Bind address is required"); }
	listen, err := NewServerSocket(web.bind);
	if err != nil { return Fmt.Errorf(
		"%s for NewServerSocket() in WebServer->Start()", err); }
	web.listen = listen;
	go web.Serve();
	Utils.SleepC();
	return nil;
}

func (web *WebServer) Serve() {
	web.service.WaitGroup.Add(1);
	defer func() {
		web.Close();
		web.service.WaitGroup.Done();
	}();
	web.service.AddClose(web);
	Log.Printf("Starting Web Server.. %s", web.bind);
	web.server.Handler = web.Router;
	if err := web.server.Serve(web.listen); err != nil {
		if !Strings.HasSuffix(err.Error(), "use of closed network connection") {
			Log.Printf("%v, in WebServer->Serve()", err); }}
}



func (web *WebServer) Close() {
	web.service.WaitGroup.Add(1);
	web.mut_state.Lock();
	defer func() {
		web.mut_state.Unlock();
		web.service.WaitGroup.Done();
	}();
	if web.listen != nil {
		if err := web.listen.Close(); err != nil {
			Log.Printf("%v, in WebServer->Close()", err); }
		if err := web.server.Shutdown(Context.Background()); err != nil {
			Log.Printf("%v, in WebServer->Close()", err); }
		web.listen = nil;
	}
}



func (web *WebServer) MiddlewareStats(next HTTP.Handler) HTTP.Handler {
	return HTTP.HandlerFunc(func(w HTTP.ResponseWriter, r *HTTP.Request) {
		start := Time.Now();
		next.ServeHTTP(w, r);
		Log.Printf("%s%s %s in %v", LogPrefixWeb, r.Method, r.URL.Path, Time.Since(start));
	});
}

func (web *WebServer) PageNotFound(w HTTP.ResponseWriter, r *HTTP.Request) {
	HTTP.Error(w, "404 Not Found", HTTP.StatusNotFound);
	Log.Printf("%s404 %s %s", LogPrefixWeb, r.Method, r.URL.Path);
}



func AddStaticRoute(router HTTP.Handler) {
	fs := HTTP.FileServer(HTTP.Dir("./static"));
	// gorilla mux
	if mux, ok := router.(*Gorilla.Router); ok {
		mux.PathPrefix("/static/").Handler(HTTP.StripPrefix("/static/", fs));
	} else
	// std http mux
	if mux, ok := router.(*HTTP.ServeMux); ok {
		mux.Handle("/static/", HTTP.StripPrefix("/static/", fs));
	// unknown mux
	} else {
		Log.Panicf("Unsupported mux type: %T", router);
	}
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
