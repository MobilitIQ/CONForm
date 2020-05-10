package main

import (
	b64 "encoding/base64"
	"fmt"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var PORT = os.Getenv("PORT")

func checkEnv() {
	if PORT == "" {
		PORT = "8812"
	}
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		if r.Method == "OPTIONS" {
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func main() {

	checkEnv()

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LFU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10 * time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
	}

	box := packr.New("Assets", "./assets")
	fileServer := http.FileServer(box)
	cached := cacheClient.Middleware(fileServer)

	http.Handle("/", CORS(cached))

	if err := http.ListenAndServe(":" + PORT, nil); err != nil {
		panic(err)
	}
}
