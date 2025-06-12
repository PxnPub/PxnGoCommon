package web;

import(
	Log      "log"
	Fmt      "fmt"
	Net      "net"
	HTTP     "net/http"
	Time     "time"
	Sync     "sync"
	Atomic   "sync/atomic"
	Gorilla  "github.com/gorilla/mux"
	Utils    "github.com/PxnPub/PxnGoCommon/utils"
	UtilsNet "github.com/PxnPub/PxnGoCommon/utils/net"
);



const LogPrefix      = "[Web] ";
const DefaultBindWeb = "tcp://127.0.0.1:8000";



type WebServer struct {
	MuxState   Sync.Mutex
	WaitGroup  *Sync.WaitGroup
	Bind       string
	Listen     Net.Listener
	Router     *Gorilla.Router
	Stats      *Stats
	NextIndex  Atomic.Uint64
	NumReqs    Atomic.Uint64
	Sessions   map[uint64]Net.Conn
}



type Stats struct {
	CountConns uint64
	CountReqs  uint64
}



func NewWebServer(bind string) *WebServer {
	web := WebServer{
		Bind:   bind,
		Router: Gorilla.NewRouter(),
	};
	web.Router.NotFoundHandler = HTTP.HandlerFunc(web.PageNotFound);
	web.Router.Use(web.MiddlewareStats);
	return &web;
}

func (web *WebServer) Start() error {
	web.MuxState.Lock();
	defer web.MuxState.Unlock();
	if web.Bind == "" { web.Bind = DefaultBindWeb; }
	listen, err := UtilsNet.NewServerSocket(web.Bind);
	if err != nil { return Fmt.Errorf(
		"%s%s for NewServerSocket in NewWebServer",
		LogPrefix, err); }
	web.Listen = listen;
	if web.WaitGroup == nil {
		var wait_group Sync.WaitGroup;
		web.WaitGroup = &wait_group;
	}
	go web.Serve();
	Utils.SleepC();
	return nil;
}

func (web *WebServer) Close() {
	web.MuxState.Lock();
	defer web.MuxState.Unlock();
	if web.Listen != nil {
		web.Listen.Close();
		web.Listen = nil;
	}
}

func (web *WebServer) CloseAll() {
	web.MuxState.Lock();
	defer web.MuxState.Unlock();
//TODO
}



func (web *WebServer) Serve() error {
	web.WaitGroup.Add(1);
	defer web.WaitGroup.Done();
	Log.Printf("Starting WebServer.. %s", web.Bind);
	return HTTP.Serve(web.Listen, web.Router);
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
