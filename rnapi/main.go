package main

import (
	"flag"
	"github.com/bitly/go-nsq"
	"github.com/julienschmidt/httprouter"
	"github.com/ymonk/readcn/pkg/alice"
	"github.com/ymonk/readcn/pkg/trace"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

var tracer trace.Tracer
var webcache *Cache
var nsqMessenger *nsq.Producer

func init() {
	tracer = trace.Log
	webcache = NewCache()
}

func main() {
	var (
		addr  = flag.String("addr", ":5050", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	flag.Parse()
	log.Println("Dialing mongo", *mongo)

	dbConn, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect mongo:", err)
	}

	defer dbConn.Close()

	// Set up NSQ
	nsqMessenger, _ = nsq.NewProducer("localhost:4150", nsq.NewConfig())
	defer nsqMessenger.Stop()

	router := newRouter()
	commonWrapper := alice.New(withLogging, withRecover, withCORS, withAPIKey, withCache, withVars, withDatabase(dbConn))

	router.GET("/articles", commonWrapper.ThenHttpRouterFunc(latestArticles))
	router.GET("/article/:permalink", commonWrapper.ThenHttpRouterFunc(readArticle))
	router.POST("/article/:permalink", commonWrapper.ThenHttpRouterFunc(editArticle))
	router.GET("/category", commonWrapper.ThenHttpRouterFunc(searchArticlesByCategory))
	router.GET("/create/id", commonWrapper.ThenHttpRouterFunc(createNewId))
	router.POST("/create/article", commonWrapper.ThenHttpRouterFunc(createArticle))
	router.GET("/delete/article/:permalink", commonWrapper.ThenHttpRouterFunc(deleteArticle))

	router.GET("/bychar", commonWrapper.ThenHttpRouterFunc(searchArticlesByChar))
	router.GET("/byvocabulary", commonWrapper.ThenHttpRouterFunc(searchArticlesByVocabulary))
	router.GET("/bygrammar", commonWrapper.ThenHttpRouterFunc(searchArticlesByGrammar))
	tracer.Trace("Starting web server on ", *addr)
	log.Fatal(http.ListenAndServe("0.0.0.0"+*addr, router))
}

// withLogging is a middleware that provides logging service
func withLogging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		tracer.Tracef("[%s] %q %v\n", r.Method, r.URL.String(), end.Sub(start))
	}
	return http.HandlerFunc(fn)
}

// withRecover is a middleware that recovers from panic
func withRecover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				tracer.Tracef("panic: %+v\n", err)
				tracer.Trace(string(debug.Stack()))
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func withAPIKey(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			http.Error(w, "invalide API key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withVars(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func withDatabase(d *mgo.Session) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			thisDB := d.Copy()
			defer thisDB.Close()
			SetDB(r, thisDB.DB("readcn_dev"))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}

func withCORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func withCache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		if !webcache.hit(url) {
			next.ServeHTTP(w, r)
			return
		}
		bytes, err := webcache.retrieve(url)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		tracer.Trace("Cache hitted, do direct retriev...")
		respondRaw(w, http.StatusOK, bytes)
	}
	return http.HandlerFunc(fn)
}

// newRouter create a new httprouter, change the NotFound handler to dull
func newRouter() (router *httprouter.Router) {
	router = httprouter.New()
	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return router
}
