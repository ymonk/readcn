package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ymonk/party"
	"github.com/ymonk/readcn/pkg/alice"
	"github.com/ymonk/readcn/pkg/trace"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var tracer trace.Tracer = trace.Log

func main() {

	dev := (AppConfig.Env == "development")

	data := map[string]interface{}{
		"ApiHostAddr": AppConfig.ApiHostAddr,
	}

	indexTemplateHandler := party.New("index.html", data, dev)
	newsTemplateHandler := party.New("news.html", data, dev)
	editTemplateHandler := party.New("edit.html", data, dev)
	editingTemplateHandler := party.New("editing.html", data, dev)
	createTemplateHandler := party.New("create.html", data, dev)
	readTemplateHandler := party.New("read.html", data, dev)
	searchTemplateHandler := party.New("search.html", data, dev)
	hskTemplateHandler := party.New("hsk.html", data, dev)
	charTemplateHandler := party.New("hsk-characters.html", data, dev)
	mucharTemplateHandler := party.New("most-used-characters.html", data, dev)
	vocabularyTemplateHandler := party.New("hsk-vocabulary.html", data, dev)
	grammarTemplateHandler := party.New("hsk-grammar.html", data, dev)
	categoryTemplateHandler := party.New("category.html", data, dev)
	exprTemplateHandler := party.New("expressions.html", data, dev)
	exprGreetingTemplateHandler := party.New("express-greetings.html", data, dev)
	exprActionTemplateHandler := party.New("express-actions.html", data, dev)
	exprOpinionTemplateHandler := party.New("express-opinions.html", data, dev)
	exprFeelingTemplateHandler := party.New("express-feelings.html", data, dev)
	exprConversationTemplateHandler := party.New("express-conversations.html", data, dev)
	eresourceTemplateHandler := party.New("eresources.html", data, dev)

	// Use httprouter as the base of the router component
	router := NewRouter()

	// commonWrapper only add logging and recover capabilities
	commonWrapper := alice.New(LoggingWrapperHandler, RecoverWrapperHandler)

	// Serve static resources
	router.Handler("GET", "/assets/*filepath",
		commonWrapper.Then(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))))

	router.Handler("GET", "/", commonWrapper.Then(indexTemplateHandler))
	router.Handler("GET", "/reading", commonWrapper.Then(newsTemplateHandler))
	router.Handler("GET", "/editing", commonWrapper.Then(editingTemplateHandler))
	router.Handler("GET", "/read", commonWrapper.Then(readTemplateHandler))
	router.Handler("GET", "/edit", commonWrapper.Then(editTemplateHandler))
	router.Handler("GET", "/create", commonWrapper.Then(createTemplateHandler))
	router.Handler("GET", "/search", commonWrapper.Then(searchTemplateHandler))
	router.Handler("POST", "/search", commonWrapper.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.FormValue("target")
		q := "c"
		if len(target) > 1 {
			if len(strings.Split(target, "/")) >= 2 {
				q = "m"
			} else {
				q = "w"
			}
		}
		tracer.Trace("Redirecting to", "/search?"+q+"="+url.QueryEscape(target))
		http.Redirect(w, r, "/search?"+q+"="+url.QueryEscape(target), http.StatusFound)
	}))
	router.Handler("GET", "/hsk", commonWrapper.Then(hskTemplateHandler))
	router.Handler("GET", "/hskchar", commonWrapper.Then(charTemplateHandler))
	router.Handler("GET", "/muchar", commonWrapper.Then(mucharTemplateHandler))
	router.Handler("GET", "/hskvocabulary", commonWrapper.Then(vocabularyTemplateHandler))
	router.Handler("GET", "/hskgrammar", commonWrapper.Then(grammarTemplateHandler))
	router.Handler("GET", "/category", commonWrapper.Then(categoryTemplateHandler))

	router.Handler("GET", "/expression", commonWrapper.Then(exprTemplateHandler))
	router.Handler("GET", "/expression/greeting", commonWrapper.Then(exprGreetingTemplateHandler))
	router.Handler("GET", "/expression/action", commonWrapper.Then(exprActionTemplateHandler))
	router.Handler("GET", "/expression/opinion", commonWrapper.Then(exprOpinionTemplateHandler))
	router.Handler("GET", "/expression/feeling", commonWrapper.Then(exprFeelingTemplateHandler))
	router.Handler("GET", "/expression/conversation", commonWrapper.Then(exprConversationTemplateHandler))

	router.Handler("GET", "/eresources", commonWrapper.Then(eresourceTemplateHandler))

	tracer.Trace("Starting web server on ", AppConfig.WebHost, AppConfig.Port)
	log.Fatal(http.ListenAndServe(AppConfig.Port, router))
}

// NewRouter create a new httprouter, change the NotFound handler to dull
func NewRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return
}
