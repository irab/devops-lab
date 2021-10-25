package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Need to have a nested struct e.g. Data{App{}}, as JSON NewEncoder.Encode returns nothing when calling an unnested struct.

type Data struct {
	App App `json:"myapplication"`
}

// App fields to pass through the pipeline
type App struct {
	Version       string `json:"version"`
	Lastcommitsha string `json:"lastcommitsha"`
	Description   string `json:"description"`
}

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	myapp := new(App)
	myapp.Version = os.Getenv("VERSION")
	myapp.Lastcommitsha = os.Getenv("LASTCOMMITSHA")
	myapp.Description = "Hello there"

	p := Data{App: *myapp}

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	r := httprouter.New()
	r.GET("/version", Version)

	log.Fatal(http.ListenAndServe(":8080", r))
}
