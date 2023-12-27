package main

import (
	"internal/dblogic"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int

	// a DB object, from the pointer of type dblogic
	DB *dblogic.DB
}

type returnVals struct {
	// the key will be the name of struct field unless you give it an explicit JSON tag
	Body string `json:"body"`
	Id   int    `json:"id"`
}

var jsonState []returnVals

// request -> http ROUTER function -> handler function
func handleRequests() {
	db, err := dblogic.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}
	//myMux := http.NewServeMux()

	myRouter := chi.NewRouter()
	// index.html is by default served, no need to specify file name
	myRouter.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	myRouter.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// use middleware to increment request counts, and write http response
	myRouter.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	myRouter.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	/*
		API logic:
		make a new chi router and group mount /healthz, /metrics, /reset all under /api
		for rApi.Get() either write a handler function explicitly, or put in function name without ()
		reqsReset is a handler that resets requests count
	*/

	rApi := chi.NewRouter()
	rApi.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})
	rApi.HandleFunc("/reset", apiCfg.reqsReset)
	rApi.HandleFunc("/chirps", apiCfg.checkInputJson)
	rApi.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	myRouter.Mount("/api", rApi)

	/*
		mount /metrics under /admin space, /admin/metrics is made this way
		if req on only root path /admin server expects to return 404
	*/
	rAdmin := chi.NewRouter()
	rAdmin.Get("/metrics", apiCfg.getAdminMetrics)
	myRouter.Mount("/admin", rAdmin)

	corsRouter := middlewareCors(myRouter)
	http.ListenAndServe(":8080", corsRouter)
}

func main() {
	dblogic.Foo()
	handleRequests()
}
