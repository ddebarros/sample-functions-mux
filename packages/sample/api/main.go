package main

import (
	"encoding/json"
	"net/http"
	"sample/api/core"
	"sample/api/gorillamux"

	"github.com/gorilla/mux"
)

var router *mux.Router
var muxAdapter *gorillamux.GorillaMuxAdapter

func init() {
	router = mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Hello world from base route")
	})

	router.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Hello world from foo route")
	})

	muxAdapter = gorillamux.New(router)
}

func Main(args map[string]interface{}) map[string]interface{} {

	mainArgs := core.MainArgsFromMap(&args)
	resp, _ := muxAdapter.MainFnAdapter(mainArgs)

	data := core.MainArgsToMap(&resp)

	if data == nil {
		m := core.ErrorResponse(500)
		return core.MainArgsToMap(&m)
	}

	return data
}
