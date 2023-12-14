package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"internal/dblogic"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

// request -> http ROUTER function -> handler function
func handleRequests() {
	apiCfg := apiConfig{0}
	//myMux := http.NewServeMux()

	myRouter := chi.NewRouter()
	// index.html is by default served, no need to specify file name
	myRouter.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	myRouter.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// use middleware to increment request counts, and write http response
	myRouter.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	myRouter.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	/*
		make a new chi router and group mount /healthz, /metrics, /reset all under /api
		for rApi.Get() either write a handler function explicitly, or put in function name without ()
		reqsReset is a handler that resets requests count
	*/
	rApi := chi.NewRouter()
	rApi.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})
	rApi.HandleFunc("/reset", apiCfg.reqsReset)
	rApi.HandleFunc("/validate_chirp", apiCfg.checkInputJson)
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

// need a handler funcition here to handle health check endpoint,
// writes OK in response message, and send status 200
func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK\n"))
	w.WriteHeader(http.StatusOK)
}

// handler function, allow cross domain access, aka other web site access this web server
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func censorWords(body string, badWords []string) string {
	words := strings.Split(body, " ")

	// Golang does not have built in In function to check if string is within a slice
	// Lane used map key checking on each input word instead of two for loops
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		for _, bw := range badWords {
			if loweredWord == bw {
				words[i] = "****"
			}
		}
	}
	cleaned_words := strings.Join(words, " ")
	return cleaned_words
}

func (cfg *apiConfig) checkInputJson(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	type returnVals struct {
		// the key will be the name of struct field unless you give it an explicit JSON tag
		Cleaned_Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	json_string := params.Body

	// check if json input is too long
	if len(json_string) > 140 {
		str := `{"error": "Something went wrong"}`
		rawJson := json.RawMessage(str)
		dat, _ := json.Marshal(rawJson)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)

	} else {
		words := []string{"fornax", "kerfuffle", "sharbert"}
		cleaned := censorWords(json_string, words)

		validJson := returnVals{
			Cleaned_Body: cleaned,
		}
		resp, err := json.Marshal(validJson)
		if err != nil {
			log.Printf("error marshalling validJson")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(resp)
	}
}

// return a handler that writes server hits back, when type apiConfig(receiver) var is calling this
// method, metod automatically have access to fields in that var
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// handler that writes number of reqests as http response
func (cfg *apiConfig) reqsCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) reqsReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
}

// handler to render admin.html
func (cfg *apiConfig) getAdminMetrics(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`

	renderedHTML := fmt.Sprintf(htmlTemplate, cfg.fileserverHits)

	// Send the rendered HTML to the HTTP response.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(renderedHTML))
}

/*
// real pointer sample code
srv := &http.Server{
    Addr:    ":" + port,
    Handler: mux,
}

log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
log.Fatal(srv.ListenAndServe())


//  Receiver of ListenAndServe is Server, like this: func (*Server) ListenAndServe
	srv is a pointer of the var of
	http.Server{
    Addr:    ":" + port,
    Handler: mux,
	}
*/

/*
 interface is a group of methods that can be given to any receiver type
 all these receiver types "satisfy" the interface adn can use any of the methods in the group
*/
