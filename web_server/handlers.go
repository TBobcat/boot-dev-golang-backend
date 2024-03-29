package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// need a handler funcition here to handle health check endpoint,
// writes OK in response message, and send status 200
func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK\n"))
	w.WriteHeader(http.StatusOK)
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

// validate input json, increment id if valid and return a json response
func (cfg *apiConfig) checkInputJson(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	// create a json file to store data (persistent state storage)
	filePath := "file_path.json"
	outputFile, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer outputFile.Close()

	// if GET method to endpoint, return jsonState
	if r.Method == http.MethodGet {
		fmt.Println("GET request received")

		//resp, err := json.Marshal(jsonState)
		resp, err := json.MarshalIndent(jsonState, "", "  ")
		if err != nil {
			log.Printf("error marshalling jsonState")
			w.WriteHeader(500)
			return
		}

		// send response json with the automatically passed in http.ResponseWriter
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(resp)
		w.Write([]byte("OK\n"))

		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		var new_id int

		// Get the last item and create id of new json
		if len(jsonState) == 0 {
			new_id = 1
		} else {
			lastIndex := len(jsonState) - 1
			new_id = jsonState[lastIndex].Id + 1
		}

		postResp := returnVals{
			Body: cleaned,
			Id:   new_id,
		}
		//resp, err := json.Marshal(postResp)
		resp, err := json.MarshalIndent(postResp, "", "  ")
		if err != nil {
			log.Printf("error marshalling postResp")
			w.WriteHeader(500)
			return
		}

		// send response json with the automatically passed in http.ResponseWriter
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(resp)

		// write updated jsonState into file
		jsonState = append(jsonState, postResp)
		updatedState, err := json.Marshal(jsonState)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		_, err = outputFile.Write(updatedState)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	// bad request ID while converting string to int.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

/*
//  pointer code example, xPtr, and &x refer to the same mem address here, so same pointer
func zero(xPtr *int) {
  *xPtr = 0
  fmt.Println(xPtr)
  fmt.Println(*xPtr)
}
func main() {
  x := 5
  zero(&x)
  fmt.Println(&x)
}


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
