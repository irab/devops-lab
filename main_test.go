package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestGetRequest(t *testing.T) {
	os.Setenv("VERSION", "1.2.3")
	os.Setenv("LASTCOMMITSHA", "8e7b64e943d5872181aaf6496d8c728edccbf696")

	router := httprouter.New()
	router.GET("/version", Version)

	req, _ := http.NewRequest("GET", "/version", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status")
	}

	// Check the response body is what we expect.
	expected := `{"myapplication":{"version":"1.2.3","lastcommitsha":"8e7b64e943d5872181aaf6496d8c728edccbf696","description":"pre-interview technical test"}}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected { // strings.TrimRight is needed as JSON Encode adds linebreaks
		t.Errorf("handler returned unexpected body: \ngot %v \nwant %v",
			rr.Body.String(), expected)
	}
}
