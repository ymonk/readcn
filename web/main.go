package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ymonk/readcn/pkg/alice"
	"github.com/ymonk/readcn/pkg/trace"
	"log"
	"net/http"
)

var tracer trace.Tracer

func init() {
	tracer = trace.Log
}

func main() {
	// Use httprouter as the base of the router component
	router := NewRouter()

	// commonWrapper only add logging and recover capabilities
	commonWrapper := alice.New(LoggingWrapperHandler, RecoverWrapperHandler)

	// Serve static resources
	router.Handler("GET", "/assets/*filepath",
		commonWrapper.Then(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))))

	// Serve static HTML
	router.Handler("GET", "/static/*filepath",
		commonWrapper.Then(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))))

	router.Handler("GET", "/", commonWrapper.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
	}))

	router.Handler("GET", "/news", commonWrapper.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/news.html", http.StatusMovedPermanently)
	}))

	router.Handler("GET", "/edit", commonWrapper.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/news.html?edit=1", http.StatusMovedPermanently)
	}))

	router.GET("/edit/:permalink", commonWrapper.ThenHttpRouterFunc(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		permalink := ps.ByName("permalink")
		http.Redirect(w, r, "/static/edit.html?v=" + permalink, http.StatusMovedPermanently)
	}))

	router.Handler("GET", "/create", commonWrapper.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/create.html", http.StatusMovedPermanently)
	}))

	tracer.Trace("Starting web server on ", AppConfig.Host, AppConfig.Port)
	log.Fatal(http.ListenAndServe(AppConfig.Host + AppConfig.Port, router))
}

// NewRouter create a new httprouter, change the NotFound handler to dull
func NewRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return
}
