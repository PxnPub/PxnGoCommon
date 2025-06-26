package net;

import(
	HTTP    "net/http"
	Gorilla "github.com/gorilla/mux"
);



type DomainsRouter struct {
	DefaultRouter HTTP.Handler
	Domains map[string]HTTP.Handler
}



func NewDomainsRouter() *DomainsRouter {
	return &DomainsRouter{
		Domains: make(map[string]HTTP.Handler),
	};
}

func (router *DomainsRouter) ServeHTTP(out HTTP.ResponseWriter, in *HTTP.Request) {
	handler, ok := router.Domains[in.Host];
	if ok { handler.ServeHTTP(out, in);
	} else if router.DefaultRouter != nil {
		router.DefaultRouter.ServeHTTP(out, in);
	} else { HTTP.NotFound(out, in); }
}



func (router *DomainsRouter) DefDomain(domain string, www bool) *Gorilla.Router {
	mux := router.AddDomain(domain, www);
	router.DefaultRouter = mux;
	return mux;
}

func (router *DomainsRouter) AddDomain(domain string, www bool) *Gorilla.Router {
	mux := Gorilla.NewRouter();
	router.Domains[domain] = mux;
	if www { router.Domains["www."+domain] = mux; }
	return mux;
}
