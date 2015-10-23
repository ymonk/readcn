package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ymonk/readcn/pkg/alice"
	"github.com/ymonk/readcn/pkg/trace"
	"log"
	"net/http"
)

var tracer trace.Tracer = trace.Log

func main() {

	dev := true

	indexTemplateHandler := NewTemplateHandler("index.html", nil, dev)
	newsTemplateHandler := NewTemplateHandler("news.html", nil, dev)
	editTemplateHandler := NewTemplateHandler("edit.html", nil, dev)
	editingTemplateHandler := NewTemplateHandler("editing.html", nil, dev)
	createTemplateHandler := NewTemplateHandler("create.html", nil, dev)
	readTemplateHandler := NewTemplateHandler("read.html", nil, dev)
	searchTemplateHandler := NewTemplateHandler("search.html", nil, dev)
	charTemplateHandler := NewTemplateHandler("hsk-characters.html", nil, dev)
	mucharTemplateHandler :=  NewTemplateHandler("most-used-characters.html", nil, dev)
	vocabularyTemplateHandler := NewTemplateHandler("hsk-vocabulary.html", nil, dev)
	grammarTemplateHandler := NewTemplateHandler("hsk-grammar.html", nil, dev)


	// Use httprouter as the base of the router component
	router := NewRouter()

	// commonWrapper only add logging and recover capabilities
	commonWrapper := alice.New(LoggingWrapperHandler, RecoverWrapperHandler)

	// Serve static resources
	router.Handler("GET", "/assets/*filepath",
		commonWrapper.Then(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))))

	// Serve static HTML
//	router.Handler("GET", "/static/*filepath",
//		commonWrapper.Then(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))))

	router.Handler("GET", "/", commonWrapper.Then(indexTemplateHandler))
	router.Handler("GET", "/reading", commonWrapper.Then(newsTemplateHandler))
	router.Handler("GET", "/editing", commonWrapper.Then(editingTemplateHandler))
	router.Handler("GET", "/read", commonWrapper.Then(readTemplateHandler))
	router.Handler("GET", "/edit", commonWrapper.Then(editTemplateHandler))
	router.Handler("GET", "/create", commonWrapper.Then(createTemplateHandler))
	router.Handler("GET", "/search", commonWrapper.Then(searchTemplateHandler))
	router.Handler("POST", "/search", commonWrapper.ThenFunc(func (w http.ResponseWriter, r *http.Request) {
		target := r.FormValue("target")
		q := "c"
		if len(target) > 1 {
			q = "w"
		}
		tracer.Trace("Redirecting to", "/search?" + q + "=" + target)
		http.Redirect(w, r, "/search?" + q + "=" + target, http.StatusFound)
	}))
	router.Handler("GET", "/hskchar", commonWrapper.Then(charTemplateHandler))
	router.Handler("GET", "/muchar", commonWrapper.Then(mucharTemplateHandler))
	router.Handler("GET", "/hskvocabulary", commonWrapper.Then(vocabularyTemplateHandler))
	router.Handler("GET", "/hskgrammar", commonWrapper.Then(grammarTemplateHandler))

//	router.GET("/edit/:permalink", commonWrapper.ThenHttpRouterFunc(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//		permalink := ps.ByName("permalink")
//		http.Redirect(w, r, "/static/edit.html?v=" + permalink, http.StatusMovedPermanently)
//	}))


	tracer.Trace("Starting web server on ", AppConfig.Host, AppConfig.Port)
	log.Fatal(http.ListenAndServe(AppConfig.Host + AppConfig.Port, router))
}

// NewRouter create a new httprouter, change the NotFound handler to dull
func NewRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return
}
